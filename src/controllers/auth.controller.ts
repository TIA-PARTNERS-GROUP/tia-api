import {
  Body,
  Controller,
  Post,
  Get,
  Route,
  SuccessResponse,
  Response,
  Tags,
  Example,
  Header,
  Security,
  Request,
} from 'tsoa';
import { StatusCodes } from 'http-status-codes';
import { ApiError } from '../errors/ApiError';
import { AuthService } from '../services/auth.service';
import { loginSchema } from '../types/auth.validation';
import type {
  LoginRequest,
  LoginResponse,
  LogoutResponse,
  SessionInfo,
  TokenValidationResponse
} from '../types/auth.dto';

/**
 * Authentication & Session Management API
 * 
 * Provides secure user authentication, session management, and token validation.
 * All endpoints use JWT tokens for authentication and maintain session integrity.
 * 
 * @security BearerAuth
 * @version 1.0.0
 */
@Route("auth")
@Tags("Authentication")
export class AuthController extends Controller {
  /**
   * User Login
   * 
   * Authenticates user credentials and creates a new session. Returns a JWT token
   * for subsequent authenticated requests. Sessions are tracked with IP and user agent.
   * 
   * @summary Authenticate user and create session
   * @description Validates user credentials and creates an authenticated session
   *              with JWT token for API access.
   * 
   * @param {LoginRequest} requestBody User login credentials
   * @param {string} user-agent User's browser/device information (auto-detected)
   * @param {string} x-forwarded-for User's IP address (auto-detected)
   * 
   * @example requestBody {
   *   "login_email": "jane.doe@example.com",
   *   "password": "SecurePassword123!"
   * }
   * 
   * @returns {LoginResponse} Authentication token and user information
   * @throws {400} Bad Request - Missing required fields
   * @throws {401} Unauthorized - Invalid credentials or inactive account
   * @throws {422} Unprocessable Entity - Validation errors
   * @throws {500} Internal Server Error - Authentication process failed
   */
  @SuccessResponse(StatusCodes.OK, "Login successful")
  @Response(StatusCodes.UNAUTHORIZED, "Invalid credentials")
  @Response(StatusCodes.UNPROCESSABLE_ENTITY, "Validation failed")
  @Example<LoginRequest>({
    login_email: "jane.doe@example.com",
    password: "SecurePassword123!"
  })
  @Post("/login")
  public async login(
    @Body() requestBody: LoginRequest,
    @Header("user-agent") userAgent?: string,
    @Header("x-forwarded-for") ipAddress?: string
  ): Promise<LoginResponse | { message: string; details?: any }> {
    try {
      const validationResult = loginSchema.safeParse(requestBody);

      if (!validationResult.success) {
        this.setStatus(StatusCodes.UNPROCESSABLE_ENTITY);
        return {
          message: 'Validation failed',
          details: validationResult.error.message,
        };
      }

      const validatedData = validationResult.data;
      const result = await AuthService.login(validatedData, ipAddress, userAgent);
      return result;
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: 'Authentication failed' };
    }
  }
  /**
   * User Logout
   * 
   * Revokes the current authentication session, invalidating the JWT token.
   * Prevents further use of the token for authenticated requests.
   * 
   * @summary Logout and revoke current session
   * @description Invalidates the current session token and logs the user out
   *              of the current device.
   * 
   * @security BearerAuth
   * 
   * @returns {LogoutResponse} Confirmation of logout
   * @throws {401} Unauthorized - Invalid or missing authentication token
   * @throws {500} Internal Server Error - Logout process failed
   */
  @SuccessResponse(StatusCodes.OK, "Logout successful")
  @Response(StatusCodes.UNAUTHORIZED, "Invalid authentication token")
  @Security("BearerAuth")
  @Post("/logout")
  public async logout(@Request() request: any): Promise<LogoutResponse | { message: string }> {
    try {
      const session = request.user?.session;
      if (!session) {
        throw new ApiError(StatusCodes.UNAUTHORIZED, 'No active session found');
      }

      const success = await AuthService.logout(session.id, session.userId);
      if (!success) {
        throw new ApiError(StatusCodes.INTERNAL_SERVER_ERROR, 'Failed to logout');
      }

      return {
        message: 'Logout successful',
        sessions_ended: 1
      };
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: 'Logout failed' };
    }
  }

  /**
   * Logout All Sessions
   * 
   * Revokes all active sessions for the authenticated user across all devices.
   * Useful for security incidents or when changing passwords.
   * 
   * @summary Logout from all devices
   * @description Invalidates all active sessions for the user across all devices
   * 
   * @security BearerAuth
   * 
   * @returns {LogoutResponse} Number of sessions ended
   * @throws {401} Unauthorized - Invalid or missing authentication token
   * @throws {500} Internal Server Error - Logout process failed
   */
  @SuccessResponse(StatusCodes.OK, "All sessions ended")
  @Response(StatusCodes.UNAUTHORIZED, "Invalid authentication token")
  @Security("BearerAuth")
  @Post("/logout-all")
  public async logoutAll(@Request() request: any): Promise<LogoutResponse | { message: string }> {
    try {
      const userId = request.user?.id;
      if (!userId) {
        throw new ApiError(StatusCodes.UNAUTHORIZED, 'User not authenticated');
      }

      const sessionsEnded = await AuthService.logoutAll(userId);

      return {
        message: 'All sessions ended successfully',
        sessions_ended: sessionsEnded
      };
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: 'Failed to logout all sessions' };
    }
  }

  /**
   * Validate Token
   * 
   * Validates the current JWT token and returns user information if valid.
   * Useful for checking token expiration and user status.
   * 
   * @summary Validate authentication token
   * @description Checks if the current JWT token is valid and returns user data
   * 
   * @security BearerAuth
   * 
   * @returns {TokenValidationResponse} Token validation result with user data
   * @throws {401} Unauthorized - Invalid or expired token
   */
  @SuccessResponse(StatusCodes.OK, "Token is valid")
  @Response(StatusCodes.UNAUTHORIZED, "Invalid or expired token")
  @Security("BearerAuth")
  @Get("/validate")
  public async validateToken(@Request() request: any): Promise<TokenValidationResponse | { message: string }> {
    try {
      const userData = request.user;

      if (!userData) {
        this.setStatus(StatusCodes.UNAUTHORIZED);
        return { message: 'Invalid token' };
      }

      const { session, ...user } = userData;

      if (!user || !session) {
        this.setStatus(StatusCodes.UNAUTHORIZED);
        return { message: 'Invalid token' };
      }

      return {
        valid: true,
        user: {
          id: user.id,
          first_name: user.first_name,
          last_name: user.last_name,
          login_email: user.login_email,
          email_verified: user.email_verified,
          active: user.active,
        },
        session: {
          id: session.id,
          ip_address: session.ip_address,
          created_at: session.created_at,
          expires_at: session.expires_at,
        }
      };
    } catch (error) {
      this.setStatus(StatusCodes.UNAUTHORIZED);
      return { message: 'Token validation failed' };
    }
  }

  /**
   * Get Active Sessions
   * 
   * Retrieves all active sessions for the authenticated user. Useful for
   * managing multiple devices and monitoring account security.
   * 
   * @summary Get user's active sessions
   * @description Returns list of all active sessions across devices
   * 
   * @security BearerAuth
   * 
   * @returns {SessionInfo[]} List of active sessions
   * @throws {401} Unauthorized - Invalid or missing authentication token
   * @throws {500} Internal Server Error - Failed to retrieve sessions
   */
  @SuccessResponse(StatusCodes.OK, "Sessions retrieved")
  @Response(StatusCodes.UNAUTHORIZED, "Invalid authentication token")
  @Security("BearerAuth")
  @Get("/sessions")
  public async getSessions(@Request() request: any): Promise<SessionInfo[] | { message: string }> {
    try {
      const userId = request.user?.id;
      if (!userId) {
        throw new ApiError(StatusCodes.UNAUTHORIZED, 'User not authenticated');
      }

      const sessions = await AuthService.getUserSessions(userId);
      return sessions;
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: 'Failed to retrieve sessions' };
    }
  }
}
