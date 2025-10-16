import { StatusCodes } from 'http-status-codes';

export class ApiError extends Error {
  public readonly statusCode: StatusCodes;
  public readonly details: any;

  constructor(statusCode: StatusCodes, message: string, details: any = null) {
    super(message);
    this.statusCode = statusCode;
    this.details = details;
  }
}
