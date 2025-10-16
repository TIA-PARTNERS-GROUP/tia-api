import jwt from 'jsonwebtoken';
import crypto from 'crypto';

export interface JWTPayload {
  userId: number;
  email: string;
  sessionId: number;
}

export class JWTUtils {
  private static readonly JWT_SECRET = process.env.JWT_SECRET || 'your-fallback-secret-change-in-production';
  private static readonly JWT_EXPIRES_IN = '7d'; // Token expires in 7 days

  /**
   * Generate JWT token for authenticated user
   */
  static generateToken(payload: Omit<JWTPayload, 'sessionId'> & { sessionId: number }): string {
    return jwt.sign(payload, this.JWT_SECRET, {
      expiresIn: this.JWT_EXPIRES_IN,
      issuer: 'tia-api',
      subject: payload.userId.toString(),
    });
  }

  /**
   * Verify and decode JWT token
   */
  static verifyToken(token: string): JWTPayload {
    try {
      return jwt.verify(token, this.JWT_SECRET) as JWTPayload;
    } catch (error) {
      if (error instanceof jwt.TokenExpiredError) {
        throw new Error('Token has expired');
      }
      if (error instanceof jwt.JsonWebTokenError) {
        throw new Error('Invalid token');
      }
      throw new Error('Token verification failed');
    }
  }

  /**
   * Hash token for secure storage in database
   */
  static hashToken(token: string): string {
    return crypto.createHash('sha256').update(token).digest('hex');
  }

  /**
   * Extract expiration date from token
   */
  static getTokenExpiry(token: string): Date {
    const decoded = jwt.decode(token) as { exp: number };
    return new Date(decoded.exp * 1000);
  }
}
