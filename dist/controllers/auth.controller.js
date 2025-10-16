"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var __param = (this && this.__param) || function (paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.AuthController = void 0;
const tsoa_1 = require("tsoa");
const http_status_codes_1 = require("http-status-codes");
const ApiError_js_1 = require("../errors/ApiError.js");
const auth_service_js_1 = require("../services/auth.service.js");
/**
 * Authentication & Session Management API
 *
 * Provides secure user authentication, session management, and token validation.
 * All endpoints use JWT tokens for authentication and maintain session integrity.
 *
 * @security BearerAuth
 * @version 1.0.0
 */
let AuthController = class AuthController extends tsoa_1.Controller {
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
    async login(requestBody, userAgent, ipAddress) {
        try {
            const result = await auth_service_js_1.AuthService.login(requestBody, ipAddress, userAgent);
            return result;
        }
        catch (error) {
            if (error instanceof ApiError_js_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
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
    async logout(request) {
        try {
            const session = request.user?.session;
            if (!session) {
                throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'No active session found');
            }
            const success = await auth_service_js_1.AuthService.logout(session.id, session.userId);
            if (!success) {
                throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR, 'Failed to logout');
            }
            return {
                message: 'Logout successful',
                sessions_ended: 1
            };
        }
        catch (error) {
            if (error instanceof ApiError_js_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
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
    async logoutAll(request) {
        try {
            const userId = request.user?.id;
            if (!userId) {
                throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'User not authenticated');
            }
            const sessionsEnded = await auth_service_js_1.AuthService.logoutAll(userId);
            return {
                message: 'All sessions ended successfully',
                sessions_ended: sessionsEnded
            };
        }
        catch (error) {
            if (error instanceof ApiError_js_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
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
    async validateToken(request) {
        try {
            const userData = request.user;
            if (!userData) {
                this.setStatus(http_status_codes_1.StatusCodes.UNAUTHORIZED);
                return { message: 'Invalid token' };
            }
            const { session, ...user } = userData;
            if (!user || !session) {
                this.setStatus(http_status_codes_1.StatusCodes.UNAUTHORIZED);
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
        }
        catch (error) {
            this.setStatus(http_status_codes_1.StatusCodes.UNAUTHORIZED);
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
    async getSessions(request) {
        try {
            const userId = request.user?.id;
            if (!userId) {
                throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'User not authenticated');
            }
            const sessions = await auth_service_js_1.AuthService.getUserSessions(userId);
            return sessions;
        }
        catch (error) {
            if (error instanceof ApiError_js_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
            return { message: 'Failed to retrieve sessions' };
        }
    }
};
exports.AuthController = AuthController;
__decorate([
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.OK, "Login successful"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNAUTHORIZED, "Invalid credentials"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNPROCESSABLE_ENTITY, "Validation failed"),
    (0, tsoa_1.Example)({
        login_email: "jane.doe@example.com",
        password: "SecurePassword123!"
    }),
    (0, tsoa_1.Post)("/login"),
    __param(0, (0, tsoa_1.Body)()),
    __param(1, (0, tsoa_1.Header)("user-agent")),
    __param(2, (0, tsoa_1.Header)("x-forwarded-for")),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object, String, String]),
    __metadata("design:returntype", Promise)
], AuthController.prototype, "login", null);
__decorate([
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.OK, "Logout successful"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNAUTHORIZED, "Invalid authentication token"),
    (0, tsoa_1.Security)("BearerAuth"),
    (0, tsoa_1.Post)("/logout"),
    __param(0, (0, tsoa_1.Request)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], AuthController.prototype, "logout", null);
__decorate([
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.OK, "All sessions ended"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNAUTHORIZED, "Invalid authentication token"),
    (0, tsoa_1.Security)("BearerAuth"),
    (0, tsoa_1.Post)("/logout-all"),
    __param(0, (0, tsoa_1.Request)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], AuthController.prototype, "logoutAll", null);
__decorate([
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.OK, "Token is valid"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNAUTHORIZED, "Invalid or expired token"),
    (0, tsoa_1.Security)("BearerAuth"),
    (0, tsoa_1.Get)("/validate"),
    __param(0, (0, tsoa_1.Request)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], AuthController.prototype, "validateToken", null);
__decorate([
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.OK, "Sessions retrieved"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNAUTHORIZED, "Invalid authentication token"),
    (0, tsoa_1.Security)("BearerAuth"),
    (0, tsoa_1.Get)("/sessions"),
    __param(0, (0, tsoa_1.Request)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], AuthController.prototype, "getSessions", null);
exports.AuthController = AuthController = __decorate([
    (0, tsoa_1.Route)("auth"),
    (0, tsoa_1.Tags)("Authentication")
], AuthController);
