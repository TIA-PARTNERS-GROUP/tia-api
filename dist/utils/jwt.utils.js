"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.JWTUtils = void 0;
const jsonwebtoken_1 = __importDefault(require("jsonwebtoken"));
const crypto_1 = __importDefault(require("crypto"));
class JWTUtils {
    static JWT_SECRET = process.env.JWT_SECRET || 'your-fallback-secret-change-in-production';
    static JWT_EXPIRES_IN = '7d'; // Token expires in 7 days
    /**
     * Generate JWT token for authenticated user
     */
    static generateToken(payload) {
        return jsonwebtoken_1.default.sign(payload, this.JWT_SECRET, {
            expiresIn: this.JWT_EXPIRES_IN,
            issuer: 'tia-api',
            subject: payload.userId.toString(),
        });
    }
    /**
     * Verify and decode JWT token
     */
    static verifyToken(token) {
        try {
            return jsonwebtoken_1.default.verify(token, this.JWT_SECRET);
        }
        catch (error) {
            if (error instanceof jsonwebtoken_1.default.TokenExpiredError) {
                throw new Error('Token has expired');
            }
            if (error instanceof jsonwebtoken_1.default.JsonWebTokenError) {
                throw new Error('Invalid token');
            }
            throw new Error('Token verification failed');
        }
    }
    /**
     * Hash token for secure storage in database
     */
    static hashToken(token) {
        return crypto_1.default.createHash('sha256').update(token).digest('hex');
    }
    /**
     * Extract expiration date from token
     */
    static getTokenExpiry(token) {
        const decoded = jsonwebtoken_1.default.decode(token);
        return new Date(decoded.exp * 1000);
    }
}
exports.JWTUtils = JWTUtils;
