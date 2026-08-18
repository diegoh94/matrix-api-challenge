export function loadEnvironment() {
  const port = Number(process.env.NODE_API_PORT ?? process.env.PORT ?? 3000);

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
