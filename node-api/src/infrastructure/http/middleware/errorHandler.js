import { DomainError } from '../../../domain/errors/DomainError.js';
import { HttpValidationError } from '../errors/HttpValidationError.js';

export function requestLogger(request, _response, next) {
  console.log(`${request.method} ${request.originalUrl}`);
  next();
}

export function errorHandler(error, _request, response, _next) {
  if (error instanceof SyntaxError && error.status === 400 && 'body' in error) {
    response.status(400).json({
      error: 'Invalid JSON payload',
      code: 'INVALID_JSON',
    });
    return;
  }

  if (error instanceof HttpValidationError || error instanceof DomainError) {
    response.status(400).json({
      error: error.message,
      code: error.code,
    });
    return;
  }

  console.error(error);

  response.status(500).json({
    error: 'Internal server error',
    code: 'INTERNAL_ERROR',
  });
}
