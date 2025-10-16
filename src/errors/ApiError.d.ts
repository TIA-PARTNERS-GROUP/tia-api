import { StatusCodes } from 'http-status-codes';
export declare class ApiError extends Error {
    readonly statusCode: StatusCodes;
    readonly details: any;
    constructor(statusCode: StatusCodes, message: string, details?: any);
}
//# sourceMappingURL=ApiError.d.ts.map