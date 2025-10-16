import { Request } from 'express';
import { AuthService } from '../services/auth.service.js';
import { ApiError } from '../errors/ApiError.js';
import { StatusCodes } from 'http-status-codes';

export async function expressAuthentication(
  request: Request,
  securityName: string,
  _scopes?: string[]
): Promise<any> {
  if (securityName === 'BearerAuth') {
    const authHeader = request.headers.authorization;
    const token = authHeader && authHeader.split(' ')[1];

    if (!token) {
      throw new ApiError(StatusCodes.UNAUTHORIZED, 'Access token required');
    }

    try {
      const { user, session } = await AuthService.validateToken(token);
      return { ...user, session };
    } catch (error) {
      throw new ApiError(StatusCodes.UNAUTHORIZED, 'Invalid authentication token');
    }
  }

  throw new ApiError(StatusCodes.UNAUTHORIZED, 'Authentication method not supported');
}
