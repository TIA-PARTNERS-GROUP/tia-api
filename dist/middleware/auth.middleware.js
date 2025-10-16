"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.expressAuthentication = expressAuthentication;
const auth_service_1 = require("../services/auth.service");
const ApiError_1 = require("../errors/ApiError");
const http_status_codes_1 = require("http-status-codes");
async function expressAuthentication(request, securityName, _scopes) {
    if (securityName === 'BearerAuth') {
        const authHeader = request.headers.authorization;
        const token = authHeader && authHeader.split(' ')[1];
        if (!token) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Access token required');
        }
        try {
            const { user, session } = await auth_service_1.AuthService.validateToken(token);
            return { ...user, session };
        }
        catch (error) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Invalid authentication token');
        }
    }
    throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Authentication method not supported');
}
