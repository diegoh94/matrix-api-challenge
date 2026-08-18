const DEFAULT_EPSILON = 1e-10;

export class DiagonalMatrixChecker {
  constructor(epsilon = DEFAULT_EPSILON) {
    this.epsilon = epsilon;
  }

  isDiagonal(matrix) {
    if (!matrix.isSquare()) {
      return false;
    }

    for (let rowIndex = 0; rowIndex < matrix.rowCount; rowIndex += 1) {
      for (let columnIndex = 0; columnIndex < matrix.columnCount; columnIndex += 1) {
        if (rowIndex === columnIndex) {
          continue;
        }

        const value = matrix.data[rowIndex][columnIndex];

        if (Math.abs(value) > this.epsilon) {
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
