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
exports.SkillsController = void 0;
const tsoa_1 = require("tsoa");
const skills_service_1 = require("../services/skills.service");
let SkillsController = class SkillsController extends tsoa_1.Controller {
    skillsService = new skills_service_1.SkillsService();
    /**
     * Get all skills with optional filtering
     */
    async getSkills(category, active, search) {
        const filters = {};
        if (category)
            filters.category = category;
        if (active !== undefined)
            filters.active = active;
        if (search)
            filters.search = search;
        return this.skillsService.getSkills(filters);
    }
    /**
     * Get a specific skill by ID
     */
    async getSkillById(skillId) {
        return this.skillsService.getSkillById(skillId);
    }
    /**
     * Get a skill by name
     */
    async getSkillByName(name) {
        return this.skillsService.getSkillByName(name);
    }
    /**
     * Create a new skill
     */
    async createSkill(body) {
        const result = await this.skillsService.createSkill(body);
        this.setStatus(201);
        return result;
    }
    /**
     * Update a skill
     */
    async updateSkill(skillId, body) {
        return this.skillsService.updateSkill(skillId, body);
    }
    /**
     * Delete a skill
     */
    async deleteSkill(skillId) {
        await this.skillsService.deleteSkill(skillId);
        this.setStatus(204);
    }
    /**
     * Get skills by category
     */
    async getSkillsByCategory(category) {
        return this.skillsService.getSkillsByCategory(category);
    }
    /**
     * Get all unique categories
     */
    async getSkillCategories() {
        return this.skillsService.getSkillCategories();
    }
    /**
     * Toggle skill active status
     */
    async toggleSkillStatus(skillId) {
        return this.skillsService.toggleSkillStatus(skillId);
    }
    /**
     * Get popular skills
     */
    async getPopularSkills(limit) {
        return this.skillsService.getPopularSkills(limit);
    }
    /**
     * Search skills
     */
    async searchSkills(query, limit) {
        return this.skillsService.searchSkills(query, limit);
    }
};
exports.SkillsController = SkillsController;
__decorate([
    (0, tsoa_1.Get)(),
    __param(0, (0, tsoa_1.Query)()),
    __param(1, (0, tsoa_1.Query)()),
    __param(2, (0, tsoa_1.Query)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Boolean, String]),
    __metadata("design:returntype", Promise)
], SkillsController.prototype, "getSkills", null);
__decorate([
    (0, tsoa_1.Get)('{skillId}'),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], SkillsController.prototype, "getSkillById", null);
__decorate([
    (0, tsoa_1.Get)('name/{name}'),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String]),
    __metadata("design:returntype", Promise)
], SkillsController.prototype, "getSkillByName", null);
__decorate([
    (0, tsoa_1.Post)(),
    (0, tsoa_1.Security)('BearerAuth'),
    (0, tsoa_1.SuccessResponse)('201', 'Created'),
    __param(0, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], SkillsController.prototype, "createSkill", null);
__decorate([
    (0, tsoa_1.Put)('{skillId}'),
    (0, tsoa_1.Security)('BearerAuth'),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, Object]),
    __metadata("design:returntype", Promise)
], SkillsController.prototype, "updateSkill", null);
__decorate([
    (0, tsoa_1.Delete)('{skillId}'),
    (0, tsoa_1.Security)('BearerAuth'),
    (0, tsoa_1.SuccessResponse)('204', 'Deleted'),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], SkillsController.prototype, "deleteSkill", null);
__decorate([
    (0, tsoa_1.Get)('category/{category}'),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String]),
    __metadata("design:returntype", Promise)
], SkillsController.prototype, "getSkillsByCategory", null);
__decorate([
    (0, tsoa_1.Get)('categories/all'),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", []),
    __metadata("design:returntype", Promise)
], SkillsController.prototype, "getSkillCategories", null);
__decorate([
    (0, tsoa_1.Put)('{skillId}/toggle-status'),
    (0, tsoa_1.Security)('BearerAuth'),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], SkillsController.prototype, "toggleSkillStatus", null);
__decorate([
    (0, tsoa_1.Get)('popular'),
    __param(0, (0, tsoa_1.Query)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], SkillsController.prototype, "getPopularSkills", null);
__decorate([
    (0, tsoa_1.Get)('search/{query}'),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Query)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Number]),
    __metadata("design:returntype", Promise)
], SkillsController.prototype, "searchSkills", null);
exports.SkillsController = SkillsController = __decorate([
    (0, tsoa_1.Route)('skills'),
    (0, tsoa_1.Tags)("Skills")
], SkillsController);
