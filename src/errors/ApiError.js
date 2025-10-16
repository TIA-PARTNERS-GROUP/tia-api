import { StatusCodes } from 'http-status-codes';
export class ApiError extends Error {
    statusCode;
    details;
    constructor(statusCode, message, details = null) {
        super(message);
        this.statusCode = statusCode;
        this.details = details;
    }
}
//# sourceMappingURL=ApiError.js.map