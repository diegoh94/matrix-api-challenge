import { DomainError, ErrorCodes } from '../../domain/errors/DomainError.js';

export class StatisticsCalculator {
  computeFromMatrices(matrices) {
    let max = Number.NEGATIVE_INFINITY;
    let min = Number.POSITIVE_INFINITY;
    let sum = 0;
    let count = 0;

    for (const matrix of matrices) {
      for (const row of matrix.data) {
        for (const value of row) {
          if (value > max) {
            max = value;
          }

          if (value < min) {
            min = value;
          }

          sum += value;
          count += 1;
        }
      }
    }

    if (count === 0) {
      throw new DomainError('No values available for statistics', ErrorCodes.NO_VALUES);
    }

    return {
      max,
      min,
      average: sum / count,
      sum,
    };
  }
}
