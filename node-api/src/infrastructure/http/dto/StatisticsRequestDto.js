import { HttpValidationError } from '../errors/HttpValidationError.js';

export class StatisticsRequestDto {
  static fromHttpBody(body) {
    if (!body || !Array.isArray(body.matrices) || body.matrices.length === 0) {
      throw new HttpValidationError('Request body must include a non-empty matrices array');
    }

    return body.matrices.map((matrix, index) => {
      if (!matrix || typeof matrix.name !== 'string' || matrix.name.trim() === '') {
        throw new HttpValidationError(`Matrix at index ${index} must include a valid name`);
      }

      if (!Array.isArray(matrix.data)) {
        throw new HttpValidationError(`Matrix at index ${index} must include a data array`);
      }

      return {
        name: matrix.name.trim(),
        data: matrix.data,
      };
    });
  }
}
