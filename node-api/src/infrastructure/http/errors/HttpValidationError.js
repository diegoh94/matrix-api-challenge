export class HttpValidationError extends Error {
  constructor(message) {
    super(message);
    this.name = 'HttpValidationError';
    this.statusCode = 400;
    this.code = 'VALIDATION_ERROR';
  }
}
