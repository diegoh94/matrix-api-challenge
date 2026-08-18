import { test } from 'node:test';
import assert from 'node:assert/strict';
import jwt from 'jsonwebtoken';

import { ComputeMatrixStatisticsUseCase } from '../../application/ComputeMatrixStatisticsUseCase.js';
import { StatisticsController } from './controllers/StatisticsController.js';
import { createServer } from './createServer.js';
import { createJwtAuthMiddleware } from './middleware/jwtAuth.js';
import { JwtTokenService } from '../auth/jwtTokenService.js';
import { StatisticsCalculator } from '../services/StatisticsCalculator.js';
import { DiagonalMatrixChecker } from '../services/DiagonalMatrixChecker.js';

const jwtSecret = 'integration-test-secret';

function buildApp({ withAuth = false } = {}) {
  const useCase = new ComputeMatrixStatisticsUseCase(
    new StatisticsCalculator(),
    new DiagonalMatrixChecker(),
  );
  const controller = new StatisticsController(useCase);
  const options = withAuth
    ? { jwtAuthMiddleware: createJwtAuthMiddleware(new JwtTokenService(jwtSecret)) }
    : {};

  return createServer(controller, options);
}

async function withServer(app, run) {
  const server = app.listen(0, '127.0.0.1');

  await new Promise((resolve) => server.once('listening', resolve));

  try {
    const { port } = server.address();
    await run(`http://127.0.0.1:${port}`);
  } finally {
    await new Promise((resolve, reject) => {
      server.close((error) => (error ? reject(error) : resolve()));
    });
  }
}

const samplePayload = {
  matrices: [
    { name: 'Q', data: [[1, 0], [0, 1]] },
    { name: 'R', data: [[2, 3], [0, 4]] },
  ],
};

test('POST /api/v1/statistics returns aggregated stats for Q and R', async () => {
  await withServer(buildApp(), async (baseUrl) => {
    const response = await fetch(`${baseUrl}/api/v1/statistics`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(samplePayload),
    });

    assert.equal(response.status, 200);

    const body = await response.json();
    assert.equal(body.max, 4);
    assert.equal(body.min, 0);
    assert.equal(body.sum, 11);
    assert.equal(body.average, 1.375);
    assert.equal(body.hasDiagonalMatrix, true);
    assert.deepEqual(body.diagonalMatrices, ['Q']);
  });
});

test('JWT middleware blocks statistics when token is missing or invalid', async () => {
  await withServer(buildApp({ withAuth: true }), async (baseUrl) => {
    const withoutToken = await fetch(`${baseUrl}/api/v1/statistics`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(samplePayload),
    });

    assert.equal(withoutToken.status, 401);

    const validToken = jwt.sign({ sub: 'matrix-api-client' }, jwtSecret, {
      algorithm: 'HS256',
      expiresIn: '1h',
    });

    const withToken = await fetch(`${baseUrl}/api/v1/statistics`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${validToken}`,
      },
      body: JSON.stringify(samplePayload),
    });

    assert.equal(withToken.status, 200);
  });
});
