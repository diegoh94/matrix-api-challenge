export class MatrixStatistics {
  constructor({ max, min, average, sum, hasDiagonalMatrix, diagonalMatrices }) {
    this.max = max;
    this.min = min;
    this.average = average;
    this.sum = sum;
    this.hasDiagonalMatrix = hasDiagonalMatrix;
    this.diagonalMatrices = diagonalMatrices;
  }
}
