import jwt from 'jsonwebtoken';

export class JwtTokenService {
  constructor(secret) {
    this.secret = secret;
  }

  validateAccessToken(accessToken) {
    try {
      jwt.verify(accessToken, this.secret, { algorithms: ['HS256'] });
      return true;
    } catch {
      return false;
    }
  }
}
