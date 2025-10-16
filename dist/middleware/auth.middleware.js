"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.expressAuthentication = expressAuthentication;
const auth_service_js_1 = require("../services/auth.service.js");
const ApiError_js_1 = require("../errors/ApiError.js");
const http_status_codes_1 = require("http-status-codes");
async function expressAuthentication(request, securityName, _scopes) {
    if (securityName === 'BearerAuth') {
        const authHeader = request.headers.authorization;
        const token = authHeader && authHeader.split(' ')[1];
        if (!token) {
            throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Access token required');
        }
        try {
            const { user, session } = await auth_service_js_1.AuthService.validateToken(token);
            return { ...user, session };
        }
        catch (error) {
            throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Invalid authentication token');
        }
    }
    throw new ApiError_js_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Authentication method not supported');
}
