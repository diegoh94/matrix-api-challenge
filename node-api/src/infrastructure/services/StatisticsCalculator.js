import { DomainError, ErrorCodes } from '../../domain/errors/DomainError.js';

export class StatisticsCalculator {
  computeFromMatrices(matrices) {
    const values = matrices.flatMap((matrix) => matrix.flattenValues());

    if (values.length === 0) {
      throw new DomainError('No values available for statistics', ErrorCodes.NO_VALUES);
    }

    const sum = values.reduce((accumulator, value) => accumulator + value, 0);

    return {
      max: Math.max(...values),
      min: Math.min(...values),
      average: sum / values.length,
      sum,
    };
  }
}
