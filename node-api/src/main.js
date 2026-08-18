import { ComputeMatrixStatisticsUseCase } from './application/ComputeMatrixStatisticsUseCase.js';
import { loadEnvironment } from './config/environment.js';
import { JwtTokenService } from './infrastructure/auth/jwtTokenService.js';
import { StatisticsController } from './infrastructure/http/controllers/StatisticsController.js';
import { createServer } from './infrastructure/http/createServer.js';
import { createJwtAuthMiddleware } from './infrastructure/http/middleware/jwtAuth.js';
import { DiagonalMatrixChecker } from './infrastructure/services/DiagonalMatrixChecker.js';
import { StatisticsCalculator } from './infrastructure/services/StatisticsCalculator.js';

const { port, jwtSecret, authEnabled } = loadEnvironment();

const statisticsCalculator = new StatisticsCalculator();
const diagonalMatrixChecker = new DiagonalMatrixChecker();
const computeMatrixStatisticsUseCase = new ComputeMatrixStatisticsUseCase(
  statisticsCalculator,
  diagonalMatrixChecker,
);
const statisticsController = new StatisticsController(computeMatrixStatisticsUseCase);

const serverOptions = authEnabled
  ? { jwtAuthMiddleware: createJwtAuthMiddleware(new JwtTokenService(jwtSecret)) }
  : {};

const app = createServer(statisticsController, serverOptions);

app.listen(port, () => {
  console.log(`node-api listening on port ${port}`);
  console.log(authEnabled ? 'JWT authentication enabled' : 'JWT authentication disabled');
});
