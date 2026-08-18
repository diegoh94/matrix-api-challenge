import express from 'express';
import { createStatisticsRoutes } from './routes/statisticsRoutes.js';
import { errorHandler, requestLogger } from './middleware/errorHandler.js';

export function createServer(statisticsController) {
  const app = express();

  app.use(express.json({ limit: '1mb' }));
  app.use(requestLogger);

  app.get('/health', statisticsController.healthCheck);
  app.use('/api/v1/statistics', createStatisticsRoutes(statisticsController));

  app.use(errorHandler);

  return app;
}
