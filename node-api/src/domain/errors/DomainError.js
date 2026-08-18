export class DomainError extends Error {
  constructor(message, code) {
    super(message);
    this.name = 'DomainError';
    this.code = code;
  }
}

export const ErrorCodes = {
  EMPTY_MATRIX: 'EMPTY_MATRIX',
  INVALID_MATRIX_SHAPE: 'INVALID_MATRIX_SHAPE',
  INVALID_MATRIX_VALUE: 'INVALID_MATRIX_VALUE',
  NO_VALUES: 'NO_VALUES',
};
