import { test } from 'node:test';
import assert from 'node:assert/strict';

import { Matrix } from '../../domain/Matrix.js';
import { DiagonalMatrixChecker } from './DiagonalMatrixChecker.js';

const checker = new DiagonalMatrixChecker();

test('identifies only square diagonal matrices by name', () => {
  const matrices = [
    new Matrix('Q', [
      [2, 0],
      [0, 3],
    ]),
    new Matrix('R', [
      [1, 2],
      [0, 4],
    ]),
    new Matrix('S', [
      [1, 2, 3],
      [4, 5, 6],
    ]),
  ];

  assert.deepEqual(checker.findDiagonalMatrixNames(matrices), ['Q']);
});

test('treats near-zero off-diagonal values as diagonal using epsilon', () => {
  const almostDiagonal = new Matrix('almost', [
    [5, 1e-12, 0],
    [0, 4, -1e-12],
    [0, 0, 2],
  ]);

  assert.equal(checker.isDiagonal(almostDiagonal), true);
});

test('rejects matrices with meaningful off-diagonal values', () => {
  const notDiagonal = new Matrix('not', [
    [5, 0.001],
    [0, 4],
  ]);

  assert.equal(checker.isDiagonal(notDiagonal), false);
});

test('rejects off-diagonal values below the diagonal', () => {
  const lowerOffDiagonal = new Matrix('lower', [
    [5, 0],
    [0.001, 4],
  ]);

  assert.equal(checker.isDiagonal(lowerOffDiagonal), false);
});
