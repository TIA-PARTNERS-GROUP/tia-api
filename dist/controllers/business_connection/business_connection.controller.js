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
exports.BusinessConnectionController = void 0;
const tsoa_1 = require("tsoa");
const http_status_codes_1 = require("http-status-codes");
const ApiError_1 = require("../../errors/ApiError");
const business_connection_service_1 = require("../../services/business_connection/business_connection.service");
let BusinessConnectionController = class BusinessConnectionController extends tsoa_1.Controller {
    connectionService = new business_connection_service_1.BusinessConnectionService();
    /**
     * Initiates a new connection request between two businesses.
     * @summary Initiate Connection
     * @returns {BusinessConnectionResponse} The new connection record with status 'pending'.
     */
    async initiateConnection(body) {
        try {
            const connection = await this.connectionService.initiateConnection(body);
            this.setStatus(http_status_codes_1.StatusCodes.CREATED);
            return connection;
        }
        catch (error) {
            if (error instanceof ApiError_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
            return { message: "Failed to initiate connection." };
        }
    }
    /**
     * Retrieves all connections associated with a specific business, with optional status filtering.
     * @summary Get Connections for Business
     * @param {number} businessId The ID of the business whose connections to retrieve.
     * @param {string} [status] Optional filter by connection status ('pending', 'active', etc.).
     * @returns {ConnectionSummaryResponse[]} A list of connection summaries.
     */
    // NOTE: This assumes the calling user has permission to view the connections of businessId.
    async getConnectionsForBusiness(businessId, status) {
        try {
            // Validation for businessId should ideally happen in a dedicated schema or middleware
            if (businessId <= 0) {
                throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.BAD_REQUEST, "Invalid Business ID.");
            }
            return await this.connectionService.getConnectionsForBusiness(businessId, status);
        }
        catch (error) {
            if (error instanceof ApiError_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
            return { message: "Failed to retrieve connections." };
        }
    }
    /**
     * Updates the status of an existing connection (e.g., accepting a pending request).
     * @summary Update Connection Status
     * @param {number} connectionId The ID of the connection record to update.
     * @returns {BusinessConnectionResponse} The updated connection record.
     */
    async updateConnectionStatus(connectionId, body) {
        try {
            return await this.connectionService.updateConnectionStatus(connectionId, body);
        }
        catch (error) {
            if (error instanceof ApiError_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
            return { message: "Failed to update connection status." };
        }
    }
    /**
     * Retrieves a single connection record by ID.
     * @summary Get Connection by ID
     * @param {number} connectionId The ID of the connection record.
     * @returns {BusinessConnectionResponse} The detailed connection record.
     */
    async getConnectionById(connectionId) {
        const connection = await this.connectionService.getConnectionById(connectionId);
        if (!connection) {
            this.setStatus(http_status_codes_1.StatusCodes.NOT_FOUND);
            return { message: "Business connection not found." };
        }
        return connection;
    }
};
exports.BusinessConnectionController = BusinessConnectionController;
__decorate([
    (0, tsoa_1.Post)(),
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.CREATED, "Connection request initiated"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.CONFLICT, "Connection already exists"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Business or User ID not found"),
    __param(0, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], BusinessConnectionController.prototype, "initiateConnection", null);
__decorate([
    (0, tsoa_1.Get)("{businessId}"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Query)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, String]),
    __metadata("design:returntype", Promise)
], BusinessConnectionController.prototype, "getConnectionsForBusiness", null);
__decorate([
    (0, tsoa_1.Put)("{connectionId}/status"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Connection not found"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.BAD_REQUEST, "Invalid status update"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object, Object]),
    __metadata("design:returntype", Promise)
], BusinessConnectionController.prototype, "updateConnectionStatus", null);
__decorate([
    (0, tsoa_1.Get)("record/{connectionId}") // Using a separate path to avoid conflict with GET /{businessId}
    ,
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Connection not found"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], BusinessConnectionController.prototype, "getConnectionById", null);
exports.BusinessConnectionController = BusinessConnectionController = __decorate([
    (0, tsoa_1.Route)("connections"),
    (0, tsoa_1.Tags)("Business Connections"),
    (0, tsoa_1.Security)("BearerAuth") // Assuming all connection endpoints require authentication
], BusinessConnectionController);
