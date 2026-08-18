import { Matrix } from '../domain/Matrix.js';
import { MatrixStatistics } from '../domain/MatrixStatistics.js';
import { DomainError, ErrorCodes } from '../domain/errors/DomainError.js';

export class ComputeMatrixStatisticsUseCase {
  constructor(statisticsCalculator, diagonalMatrixChecker) {
    this.statisticsCalculator = statisticsCalculator;
    this.diagonalMatrixChecker = diagonalMatrixChecker;
  }

  execute(matrixInputs) {
    if (!Array.isArray(matrixInputs) || matrixInputs.length === 0) {
      throw new DomainError('At least one matrix is required', ErrorCodes.EMPTY_MATRIX);
    }

    const matrices = matrixInputs.map(({ name, data }) => new Matrix(name, data));

    const numericStatistics = this.statisticsCalculator.computeFromMatrices(matrices);
    const diagonalMatrices = this.diagonalMatrixChecker.findDiagonalMatrixNames(matrices);

    return new MatrixStatistics({
      ...numericStatistics,
      hasDiagonalMatrix: diagonalMatrices.length > 0,
      diagonalMatrices,
    });
  }
}
