export function createJwtAuthMiddleware(jwtTokenService) {
  return function jwtAuth(request, response, next) {
    const authorizationHeader = request.headers.authorization ?? '';

    if (!authorizationHeader.startsWith('Bearer ')) {
      response.status(401).json({
        error: 'missing or invalid authorization header',
        code: 'UNAUTHORIZED',
      });
      return;
    }

    const accessToken = authorizationHeader.slice('Bearer '.length).trim();

    if (!accessToken || !jwtTokenService.validateAccessToken(accessToken)) {
      response.status(401).json({
        error: 'invalid or expired token',
        code: 'UNAUTHORIZED',
      });
      return;
    }

    next();
  };
}
