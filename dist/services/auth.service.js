"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AuthService = void 0;
const prisma_js_1 = require("../lib/prisma.js");
const http_status_codes_1 = require("http-status-codes");
const ApiError_js_1 = require("../errors/ApiError.js");
const jwt_utils_js_1 = require("../utils/jwt.utils.js");
const password_utils_js_1 = require("../utils/password.utils.js");
class AuthService {
    /**
     * Authenticate user and create session
     */
    static async login(loginData, ipAddress, userAgent) {
        const { login_email, password } = loginData;
        const user = await prisma_js_1.prisma.users.findUnique({
            where: { login_email },
            select: {
                id: true,
                first_name: true,
                last_name: true,
                login_email: true,
                password_hash: true,
                contact_email: true,
                email_verified: true,
                active: true,
                created_at: true,
            }
        });
        if (!user) {
            throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Invalid email or password');
        }
        if (!user.active) {
            throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Account is deactivated');
        }
        if (!user.password_hash) {
            throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'No password set for this account');
        }
        const isPasswordValid = await password_utils_js_1.PasswordUtils.verifyPassword(password, user.password_hash);
        if (!isPasswordValid) {
            throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Invalid email or password');
        }
        const placeholderToken = "pending_" + Date.now();
        const sessionData = {
            user_id: user.id,
            token_hash: jwt_utils_js_1.JWTUtils.hashToken(placeholderToken), // Temporary hash
            expires_at: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000), // 7 days from now
        };
        if (ipAddress)
            sessionData.ip_address = ipAddress;
        if (userAgent)
            sessionData.user_agent = userAgent;
        const session = await prisma_js_1.prisma.user_sessions.create({
            data: sessionData
        });
        const realToken = jwt_utils_js_1.JWTUtils.generateToken({
            userId: user.id,
            email: user.login_email,
            sessionId: session.id,
        });
        await prisma_js_1.prisma.user_sessions.update({
            where: { id: session.id },
            data: {
                token_hash: jwt_utils_js_1.JWTUtils.hashToken(realToken),
                expires_at: jwt_utils_js_1.JWTUtils.getTokenExpiry(realToken)
            }
        });
        const { password_hash, ...userWithoutPassword } = user;
        return {
            user: userWithoutPassword,
            token: realToken,
            session_id: session.id,
            expires_at: session.expires_at,
            token_type: 'Bearer'
        };
    }
    /**
     * Logout user by revoking session
     */
    static async logout(sessionId, userId) {
        const session = await prisma_js_1.prisma.user_sessions.updateMany({
            where: {
                id: sessionId,
                user_id: userId,
                revoked_at: null
            },
            data: {
                revoked_at: new Date()
            }
        });
        return session.count > 0;
    }
    /**
     * Logout all sessions for user
     */
    static async logoutAll(userId) {
        const result = await prisma_js_1.prisma.user_sessions.updateMany({
            where: {
                user_id: userId,
                revoked_at: null,
                expires_at: { gt: new Date() } // Only active sessions
            },
            data: {
                revoked_at: new Date()
            }
        });
        return result.count;
    }
    /**
     * Validate JWT token and return user
     */
    static async validateToken(token) {
        try {
            jwt_utils_js_1.JWTUtils.verifyToken(token);
            const tokenHash = jwt_utils_js_1.JWTUtils.hashToken(token);
            const session = await prisma_js_1.prisma.user_sessions.findFirst({
                where: {
                    token_hash: tokenHash,
                    revoked_at: null,
                    expires_at: { gt: new Date() }
                },
                include: {
                    users: {
                        select: {
                            id: true,
                            first_name: true,
                            last_name: true,
                            login_email: true,
                            contact_email: true,
                            email_verified: true,
                            active: true,
                        }
                    }
                }
            });
            if (!session) {
                throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Invalid or expired session');
            }
            if (!session.users.active) {
                throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Account is deactivated');
            }
            return {
                user: session.users,
                session: {
                    id: session.id,
                    ip_address: session.ip_address,
                    created_at: session.created_at,
                    expires_at: session.expires_at
                }
            };
        }
        catch (error) {
            if (error instanceof ApiError_js_1.ApiError) {
                throw error;
            }
            throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Invalid authentication token');
        }
    }
    /**
     * Get active sessions for user
     */
    static async getUserSessions(userId) {
        const sessions = await prisma_js_1.prisma.user_sessions.findMany({
            where: {
                user_id: userId,
                revoked_at: null,
                expires_at: { gt: new Date() }
            },
            select: {
                id: true,
                ip_address: true,
                user_agent: true,
                created_at: true,
                expires_at: true,
            },
            orderBy: {
                created_at: 'desc'
            }
        });
        return sessions;
    }
    /**
     * Clean up expired sessions
     */
    static async cleanupExpiredSessions() {
        const result = await prisma_js_1.prisma.user_sessions.deleteMany({
            where: {
                OR: [
                    { expires_at: { lt: new Date() } },
                    { revoked_at: { not: null } }
                ]
            }
        });
        return result.count;
    }
}
exports.AuthService = AuthService;
