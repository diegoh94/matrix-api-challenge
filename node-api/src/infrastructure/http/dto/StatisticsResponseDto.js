const JSON_DECIMAL_PLACES = 6;

function roundToDecimalPlaces(value, decimalPlaces) {
  const factor = 10 ** decimalPlaces;
  return Math.round(value * factor) / factor;
}

export class StatisticsResponseDto {
  static fromDomain(statistics) {
    return {
      max: roundToDecimalPlaces(statistics.max, JSON_DECIMAL_PLACES),
      min: roundToDecimalPlaces(statistics.min, JSON_DECIMAL_PLACES),
      average: roundToDecimalPlaces(statistics.average, JSON_DECIMAL_PLACES),
      sum: roundToDecimalPlaces(statistics.sum, JSON_DECIMAL_PLACES),
      hasDiagonalMatrix: statistics.hasDiagonalMatrix,
      diagonalMatrices: statistics.diagonalMatrices,
    };
  }
}
