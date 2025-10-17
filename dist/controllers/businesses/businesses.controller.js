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
exports.BusinessesController = void 0;
const tsoa_1 = require("tsoa");
const businesses_service_1 = require("../../services/businesses/businesses.service");
const base_controller_1 = require("../../controllers/base.controller");
const http_status_codes_1 = require("http-status-codes");
/**
 * Business Management API
 *
 * Provides endpoints for creating, retrieving, updating, and deleting business profiles.
 * Also includes endpoints for managing business status and retrieving statistics.
 *
 * @version 1.0.0
 */
let BusinessesController = class BusinessesController extends base_controller_1.BaseController {
    businessesService = new businesses_service_1.BusinessesService();
    /**
     * Get All Businesses with Filtering
     *
     * Retrieves a list of all businesses, with optional query parameters for filtering
     * by type, category, phase, status, and operator user ID.
     *
     * @summary Get a filtered list of businesses
     * @returns {BusinessResponse[]} A list of business profiles.
     */
    async getBusinesses(business_type, business_category, business_phase, active, search, operator_user_id) {
        try {
            const filters = {
                business_type,
                business_category,
                business_phase,
                active,
                search,
                operator_user_id,
            };
            return await this.businessesService.getBusinesses(filters);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Get Business by ID
     *
     * Retrieves detailed information for a specific business by its unique identifier.
     *
     * @summary Get a single business by its ID
     * @param {number} businessId The unique identifier of the business.
     * @returns {BusinessResponse} The requested business profile.
     * @throws {404} Not Found - The business with the specified ID does not exist.
     */
    async getBusinessById(businessId) {
        try {
            return await this.businessesService.getBusinessById(businessId);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Create a New Business
     *
     * Registers a new business profile in the system. Requires authentication.
     *
     * @summary Create a new business profile
     * @param {CreateBusinessInput} body The data for the new business.
     * @returns {BusinessResponse} The newly created business profile.
     * @throws {401} Unauthorized - Authentication is required.
     * @throws {422} Unprocessable Entity - Validation errors in input data.
     */
    async createBusiness(body) {
        try {
            const result = await this.businessesService.createBusiness(body);
            this.setStatus(http_status_codes_1.StatusCodes.CREATED);
            return result;
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Update a Business
     *
     * Modifies the details of an existing business profile. Requires authentication.
     *
     * @summary Update an existing business
     * @param {number} businessId The unique identifier of the business to update.
     * @param {UpdateBusinessInput} body The fields to update.
     * @returns {BusinessResponse} The updated business profile.
     * @throws {401} Unauthorized - Authentication is required.
     * @throws {404} Not Found - The business with the specified ID does not exist.
     */
    async updateBusiness(businessId, body) {
        try {
            return await this.businessesService.updateBusiness(businessId, body);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Delete a Business
     *
     * Permanently removes a business profile from the system. Requires authentication.
     *
     * @summary Delete a business profile
     * @param {number} businessId The unique identifier of the business to delete.
     * @throws {204} No Content - The business was successfully deleted.
     * @throws {401} Unauthorized - Authentication is required.
     * @throws {404} Not Found - The business with the specified ID does not exist.
     */
    async deleteBusiness(businessId) {
        try {
            await this.businessesService.deleteBusiness(businessId);
            this.setStatus(http_status_codes_1.StatusCodes.NO_CONTENT);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Get All Businesses for a User
     *
     * Retrieves a list of all businesses associated with a specific user.
     *
     * @summary Get all businesses for a user
     * @param {number} userId The unique identifier of the user.
     * @returns {BusinessResponse[]} A list of business profiles.
     */
    async getUserBusinesses(userId) {
        try {
            return await this.businessesService.getUserBusinesses(userId);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Toggle Business Active Status
     *
     * Toggles the active status of a business. Requires authentication.
     *
     * @summary Toggle the active status of a business
     * @param {number} businessId The unique identifier of the business.
     * @returns {BusinessResponse} The business profile with the updated status.
     * @throws {401} Unauthorized - Authentication is required.
     * @throws {404} Not Found - The business with the specified ID does not exist.
     */
    async toggleBusinessStatus(businessId) {
        try {
            return await this.businessesService.toggleBusinessStatus(businessId);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Get Business Statistics
     *
     * Retrieves statistics and analytics for a specific business.
     *
     * @summary Get statistics for a business
     * @param {number} businessId The unique identifier of the business.
     * @returns {BusinessStatsResponse} The statistics for the business.
     * @throws {404} Not Found - The business with the specified ID does not exist.
     */
    async getBusinessStats(businessId) {
        try {
            return await this.businessesService.getBusinessStats(businessId);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
};
exports.BusinessesController = BusinessesController;
__decorate([
    (0, tsoa_1.Get)(),
    __param(0, (0, tsoa_1.Query)()),
    __param(1, (0, tsoa_1.Query)()),
    __param(2, (0, tsoa_1.Query)()),
    __param(3, (0, tsoa_1.Query)()),
    __param(4, (0, tsoa_1.Query)()),
    __param(5, (0, tsoa_1.Query)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, String, Boolean, String, Number]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "getBusinesses", null);
__decorate([
    (0, tsoa_1.Get)("{businessId}"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Business not found"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "getBusinessById", null);
__decorate([
    (0, tsoa_1.Post)(),
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.CREATED, http_status_codes_1.ReasonPhrases.CREATED),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNAUTHORIZED, http_status_codes_1.ReasonPhrases.UNAUTHORIZED),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNPROCESSABLE_ENTITY, http_status_codes_1.ReasonPhrases.UNPROCESSABLE_ENTITY),
    __param(0, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "createBusiness", null);
__decorate([
    (0, tsoa_1.Put)("{businessId}"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNAUTHORIZED, http_status_codes_1.ReasonPhrases.UNAUTHORIZED),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Business not found"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, Object]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "updateBusiness", null);
__decorate([
    (0, tsoa_1.Delete)("{businessId}"),
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.NO_CONTENT, "Deleted"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNAUTHORIZED, http_status_codes_1.ReasonPhrases.UNAUTHORIZED),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Business not found"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "deleteBusiness", null);
__decorate([
    (0, tsoa_1.Get)("user/{userId}"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "getUserBusinesses", null);
__decorate([
    (0, tsoa_1.Put)("{businessId}/toggle-status"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNAUTHORIZED, http_status_codes_1.ReasonPhrases.UNAUTHORIZED),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Business not found"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "toggleBusinessStatus", null);
__decorate([
    (0, tsoa_1.Get)("{businessId}/stats"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Business not found"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "getBusinessStats", null);
exports.BusinessesController = BusinessesController = __decorate([
    (0, tsoa_1.Route)("businesses"),
    (0, tsoa_1.Tags)("Businesses"),
    (0, tsoa_1.Security)("BearerAuth"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR, http_status_codes_1.ReasonPhrases.INTERNAL_SERVER_ERROR)
], BusinessesController);
