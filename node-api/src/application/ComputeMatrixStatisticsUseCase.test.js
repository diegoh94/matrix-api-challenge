import { test } from 'node:test';
import assert from 'node:assert/strict';

import { ComputeMatrixStatisticsUseCase } from './ComputeMatrixStatisticsUseCase.js';
import { StatisticsCalculator } from '../infrastructure/services/StatisticsCalculator.js';
import { DiagonalMatrixChecker } from '../infrastructure/services/DiagonalMatrixChecker.js';
import { DomainError, ErrorCodes } from '../domain/errors/DomainError.js';

const useCase = new ComputeMatrixStatisticsUseCase(
  new StatisticsCalculator(),
  new DiagonalMatrixChecker(),
);

test('aggregates statistics across Q and R like the production payload', () => {
  const result = useCase.execute([
    { name: 'Q', data: [[1, 0, 0], [0, 1, 0], [0, 0, 1]] },
    { name: 'R', data: [[-5.91608, -7.437357], [0, 0.828079], [0, 0]] },
  ]);

  assert.equal(result.max, 1);
  assert.equal(result.min, -7.437357);
  assert.ok(result.sum < 0);
  assert.ok(result.average < 0);
  assert.equal(result.hasDiagonalMatrix, true);
  assert.deepEqual(result.diagonalMatrices, ['Q']);
});

test('rejects empty matrix list before touching calculators', () => {
  assert.throws(
    () => useCase.execute([]),
    (error) => error instanceof DomainError && error.code === ErrorCodes.EMPTY_MATRIX,
  );
});

test('rejects invalid matrix shape from domain validation', () => {
  assert.throws(
    () =>
      useCase.execute([
        { name: 'Q', data: [[1, 2], [3]] },
        { name: 'R', data: [[4, 5]] },
      ]),
    (error) => error instanceof DomainError && error.code === ErrorCodes.INVALID_MATRIX_SHAPE,
  );
});
