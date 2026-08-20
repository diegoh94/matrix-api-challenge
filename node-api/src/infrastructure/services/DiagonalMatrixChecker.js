const DEFAULT_EPSILON = 1e-10;

export class DiagonalMatrixChecker {
  constructor(epsilon = DEFAULT_EPSILON) {
    this.epsilon = epsilon;
  }

  isDiagonal(matrix) {
    if (!matrix.isSquare()) {
      return false;
    }

    const size = matrix.rowCount;

    for (let rowIndex = 0; rowIndex < size; rowIndex += 1) {
      for (let columnIndex = rowIndex + 1; columnIndex < size; columnIndex += 1) {
        if (
          Math.abs(matrix.data[rowIndex][columnIndex]) > this.epsilon ||
          Math.abs(matrix.data[columnIndex][rowIndex]) > this.epsilon
        ) {
          return false;
        }
      }
    }

    return true;
  }

  findDiagonalMatrixNames(matrices) {
    return matrices.filter((matrix) => this.isDiagonal(matrix)).map((matrix) => matrix.name);
  }
}
