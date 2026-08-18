import { Router } from 'express';

export function createStatisticsRoutes(statisticsController, jwtAuthMiddleware = null) {
  const router = Router();

  if (jwtAuthMiddleware) {
    router.post('/', jwtAuthMiddleware, statisticsController.computeStatistics);
  } else {
    router.post('/', statisticsController.computeStatistics);
  }

  return router;
}
