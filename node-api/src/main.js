import { ComputeMatrixStatisticsUseCase } from './application/ComputeMatrixStatisticsUseCase.js';
import { loadEnvironment } from './config/environment.js';
import { StatisticsController } from './infrastructure/http/controllers/StatisticsController.js';
import { createServer } from './infrastructure/http/createServer.js';
import { DiagonalMatrixChecker } from './infrastructure/services/DiagonalMatrixChecker.js';
import { StatisticsCalculator } from './infrastructure/services/StatisticsCalculator.js';

const { port } = loadEnvironment();

const statisticsCalculator = new StatisticsCalculator();
const diagonalMatrixChecker = new DiagonalMatrixChecker();
const computeMatrixStatisticsUseCase = new ComputeMatrixStatisticsUseCase(
  statisticsCalculator,
  diagonalMatrixChecker,
);
const statisticsController = new StatisticsController(computeMatrixStatisticsUseCase);
const app = createServer(statisticsController);

app.listen(port, () => {
  console.log(`node-api listening on port ${port}`);
});
