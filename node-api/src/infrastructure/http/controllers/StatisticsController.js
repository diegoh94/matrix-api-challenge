import { StatisticsRequestDto } from '../dto/StatisticsRequestDto.js';
import { StatisticsResponseDto } from '../dto/StatisticsResponseDto.js';

export class StatisticsController {
  constructor(computeMatrixStatisticsUseCase) {
    this.computeMatrixStatisticsUseCase = computeMatrixStatisticsUseCase;
  }

  computeStatistics = (request, response, next) => {
    try {
      const matrixInputs = StatisticsRequestDto.fromHttpBody(request.body);
      const statistics = this.computeMatrixStatisticsUseCase.execute(matrixInputs);

      response.status(200).json(StatisticsResponseDto.fromDomain(statistics));
    } catch (error) {
      next(error);
    }
  };

  healthCheck = (_request, response) => {
    response.status(200).json({
      status: 'ok',
      service: 'node-api',
    });
  };
}
