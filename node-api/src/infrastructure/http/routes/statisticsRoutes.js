import { Router } from 'express';

export function createStatisticsRoutes(statisticsController) {
  const router = Router();

  router.post('/', statisticsController.computeStatistics);

  return router;
}
