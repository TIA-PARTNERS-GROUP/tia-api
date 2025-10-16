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
const businesses_service_1 = require("../services/businesses.service");
let BusinessesController = class BusinessesController extends tsoa_1.Controller {
    businessesService = new businesses_service_1.BusinessesService();
    async getBusinesses(business_type, business_category, business_phase, active, search, operator_user_id) {
        const filters = {};
        if (business_type)
            filters.business_type = business_type;
        if (business_category)
            filters.business_category = business_category;
        if (business_phase)
            filters.business_phase = business_phase;
        if (active !== undefined)
            filters.active = active;
        if (search)
            filters.search = search;
        if (operator_user_id)
            filters.operator_user_id = operator_user_id;
        return this.businessesService.getBusinesses(filters);
    }
    async getBusinessById(businessId) {
        return this.businessesService.getBusinessById(businessId);
    }
    async createBusiness(body) {
        const result = await this.businessesService.createBusiness(body);
        this.setStatus(201);
        return result;
    }
    async updateBusiness(businessId, body) {
        return this.businessesService.updateBusiness(businessId, body);
    }
    async deleteBusiness(businessId) {
        await this.businessesService.deleteBusiness(businessId);
        this.setStatus(204);
    }
    async getUserBusinesses(userId) {
        return this.businessesService.getUserBusinesses(userId);
    }
    async toggleBusinessStatus(businessId) {
        return this.businessesService.toggleBusinessStatus(businessId);
    }
    async getBusinessStats(businessId) {
        return this.businessesService.getBusinessStats(businessId);
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
    (0, tsoa_1.Get)('{businessId}'),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "getBusinessById", null);
__decorate([
    (0, tsoa_1.Post)(),
    (0, tsoa_1.Security)("BearerAuth"),
    (0, tsoa_1.SuccessResponse)('201', 'Created'),
    __param(0, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "createBusiness", null);
__decorate([
    (0, tsoa_1.Put)('{businessId}'),
    (0, tsoa_1.Security)("BearerAuth"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, Object]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "updateBusiness", null);
__decorate([
    (0, tsoa_1.Delete)('{businessId}'),
    (0, tsoa_1.Security)("BearerAuth"),
    (0, tsoa_1.SuccessResponse)('204', 'Deleted'),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "deleteBusiness", null);
__decorate([
    (0, tsoa_1.Get)('user/{userId}'),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "getUserBusinesses", null);
__decorate([
    (0, tsoa_1.Put)('{businessId}/toggle-status'),
    (0, tsoa_1.Security)("BearerAuth"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "toggleBusinessStatus", null);
__decorate([
    (0, tsoa_1.Get)('{businessId}/stats'),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], BusinessesController.prototype, "getBusinessStats", null);
exports.BusinessesController = BusinessesController = __decorate([
    (0, tsoa_1.Route)('businesses')
], BusinessesController);
