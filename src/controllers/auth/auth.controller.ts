import {
  Body,
  Post,
  Get,
  Route,
  SuccessResponse,
  Response,
  Tags,
  Header,
  Security,
  Request,
} from "tsoa";
import type {
  LoginResponse,
  LogoutResponse,
  SessionInfo,
  TokenValidationResponse,
} from "types/auth/auth.dto.js";
import type { LoginInput } from "types/auth/auth.validation.js";
import { StatusCodes } from "http-status-codes";
import { HttpErrors, ApiErrorResponse } from "errors/ApiError.js";
import { AuthService } from "services/auth/auth.service.js";
import { BaseController } from "controllers/base.controller.js";

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
export class AuthController extends BaseController {
  /**
   * User Login
   *
   * Authenticates user credentials and creates a new session. Returns a JWT token
   * for subsequent authenticated requests. Sessions are tracked with IP and user agent.
   *
   * @summary Authenticate user and create session
   * @description Validates user credentials and creates an authenticated session
   * with JWT token for API access.
   *
   * @param {LoginInput} requestBody User login credentials
   * @param {string} user-agent User's browser/device information (auto-detected)
   * @param {string} x-forwarded-for User's IP address (auto-detected)
   *
   * @returns {LoginResponse} Authentication token and user information
   * @throws {400} Bad Request - Missing required fields
   * @throws {401} Unauthorized - Invalid credentials or inactive account
   * @throws {422} Unprocessable Entity - Validation errors
   * @throws {500} Internal Server Error - Authentication process failed
   */
  @SuccessResponse(StatusCodes.OK, "Login successful")
  @Response<ApiErrorResponse>(StatusCodes.UNAUTHORIZED, "Invalid credentials")
  @Response<ApiErrorResponse>(
    StatusCodes.UNPROCESSABLE_ENTITY,
    "Validation failed",
  )
  @Post("/login")
  public async login(
    @Body() requestBody: LoginInput,
    @Header("user-agent") userAgent?: string,
    @Header("x-forwarded-for") ipAddress?: string,
  ): Promise<LoginResponse | ApiErrorResponse> {
    try {
      return await AuthService.login(requestBody, ipAddress, userAgent);
    } catch (error) {
      return this.handleError(error);
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
   * of the current device.
   *
   * @security BearerAuth
   *
   * @returns {LogoutResponse} Confirmation of logout
   * @throws {401} Unauthorized - Invalid or missing authentication token
   * @throws {500} Internal Server Error - Logout process failed
   */
  @SuccessResponse(StatusCodes.OK, "Logout successful")
  @Response<ApiErrorResponse>(
    StatusCodes.UNAUTHORIZED,
    "Invalid authentication token",
  )
  @Security("BearerAuth")
  @Post("/logout")
  public async logout(
    @Request() request: any,
  ): Promise<LogoutResponse | ApiErrorResponse> {
    try {
      const sessionId = request.user?.session?.id;
      const userId = request.user?.id;

      if (!sessionId || !userId) {
        throw HttpErrors.Unauthorized("No active session found");
      }

      await AuthService.logout(sessionId, userId);

      return { message: "Logout successful", sessions_ended: 1 };
    } catch (error) {
      return this.handleError(error);
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
  @Response<ApiErrorResponse>(
    StatusCodes.UNAUTHORIZED,
    "Invalid authentication token",
  )
  @Security("BearerAuth")
  @Post("/logout-all")
  public async logoutAll(
    @Request() request: any,
  ): Promise<LogoutResponse | ApiErrorResponse> {
    try {
      const userId = request.user?.id;
      if (!userId) {
        throw HttpErrors.Unauthorized("User not authenticated");
      }
      const sessionsEnded = await AuthService.logoutAll(userId);
      return {
        message: "All sessions ended successfully",
        sessions_ended: sessionsEnded,
      };
    } catch (error) {
      return this.handleError(error);
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
  @Response<ApiErrorResponse>(
    StatusCodes.UNAUTHORIZED,
    "Invalid or expired token",
  )
  @Security("BearerAuth")
  @Get("/validate")
  public async validateToken(
    @Request() request: any,
  ): Promise<TokenValidationResponse | ApiErrorResponse> {
    try {
      const { session, ...user } = request.user || {};

      if (!user?.id || !session?.id) {
        throw HttpErrors.Unauthorized(
          "Invalid token: user or session data missing",
        );
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
        },
      };
    } catch (error) {
      return this.handleError(error);
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
  @Response<ApiErrorResponse>(
    StatusCodes.UNAUTHORIZED,
    "Invalid authentication token",
  )
  @Security("BearerAuth")
  @Get("/sessions")
  public async getSessions(
    @Request() request: any,
  ): Promise<SessionInfo[] | ApiErrorResponse> {
    try {
      const userId = request.user?.id;
      if (!userId) {
        throw HttpErrors.Unauthorized("User not authenticated");
      }
      return await AuthService.getUserSessions(userId);
    } catch (error) {
      return this.handleError(error);
    }
  }
}
