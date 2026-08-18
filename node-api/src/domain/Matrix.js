import { DomainError, ErrorCodes } from './errors/DomainError.js';

export class Matrix {
  constructor(name, data) {
    this.name = name;
    this.data = data;
    this.#validate();
  }

  get rowCount() {
    return this.data.length;
  }

  get columnCount() {
    return this.data[0]?.length ?? 0;
  }

  isSquare() {
    return this.rowCount === this.columnCount;
  }

  flattenValues() {
    return this.data.flat();
  }

  #validate() {
    if (!Array.isArray(this.data) || this.data.length === 0) {
      throw new DomainError('Matrix must contain at least one row', ErrorCodes.EMPTY_MATRIX);
    }

    const expectedColumnCount = this.data[0].length;

    if (expectedColumnCount === 0) {
      throw new DomainError('Matrix rows must contain at least one column', ErrorCodes.EMPTY_MATRIX);
    }

    for (const row of this.data) {
      if (!Array.isArray(row) || row.length !== expectedColumnCount) {
        throw new DomainError(
          'All matrix rows must have the same number of columns',
          ErrorCodes.INVALID_MATRIX_SHAPE,
        );
      }

      for (const value of row) {
        if (typeof value !== 'number' || !Number.isFinite(value)) {
          throw new DomainError(
            'Matrix values must be finite numbers',
            ErrorCodes.INVALID_MATRIX_VALUE,
          );
        }
      }
    }
  }
}
