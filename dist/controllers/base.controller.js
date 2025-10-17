"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BaseController = void 0;
const tsoa_1 = require("tsoa");
const ApiError_1 = require("../errors/ApiError");
class BaseController extends tsoa_1.Controller {
    /**
     * A consistent error handling method for all controllers.
     * It checks if the error is a known ApiError and formats it correctly.
     * If it's an unknown error, it logs it and returns a generic 500 response,
     * hiding implementation details in production.
     * @param error The error caught in a try/catch block.
     * @returns A formatted ApiErrorResponse.
     */
    handleError(error) {
        if (error instanceof ApiError_1.ApiError) {
            this.setStatus(error.statusCode);
            return error.toResponse();
        }
        console.error("An unexpected error occurred:", error);
        const internalError = ApiError_1.HttpErrors.InternalServerError();
        this.setStatus(internalError.statusCode);
        if (process.env.NODE_ENV === "production") {
            return {
                message: internalError.message,
                statusCode: internalError.statusCode,
            };
        }
        return {
            message: internalError.message,
            statusCode: internalError.statusCode,
            details: error instanceof Error ? error.stack : String(error),
        };
    }
}
exports.BaseController = BaseController;
