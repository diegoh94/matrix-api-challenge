export function loadEnvironment() {
  // PORT first: Railway/cloud injects it; NODE_API_PORT for Docker Compose local.
  const port = Number(process.env.PORT ?? process.env.NODE_API_PORT ?? 3000);

  if (!Number.isInteger(port) || port <= 0) {
    throw new Error('NODE_API_PORT must be a positive integer');
  }

  const jwtSecret = process.env.JWT_SECRET ?? '';

  return {
    port,
    jwtSecret,
    authEnabled: jwtSecret.length > 0,
  };
}
