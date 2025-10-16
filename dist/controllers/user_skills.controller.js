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
exports.UserSkillsController = void 0;
const tsoa_1 = require("tsoa");
const user_skills_service_1 = require("../services/user_skills.service");
let UserSkillsController = class UserSkillsController extends tsoa_1.Controller {
    userSkillsService = new user_skills_service_1.UserSkillsService();
    /**
     * Get all skills for a user
     */
    async getUserSkills(userId) {
        return this.userSkillsService.getUserSkills(userId);
    }
    /**
     * Get a specific user skill
     */
    async getUserSkillById(skillId, userId) {
        return this.userSkillsService.getUserSkillById(skillId, userId);
    }
    /**
     * Create a new user skill
     */
    async createUserSkill(body) {
        const result = await this.userSkillsService.createUserSkill(body);
        this.setStatus(201);
        return result;
    }
    /**
     * Update a user skill
     */
    async updateUserSkill(skillId, userId, body) {
        return this.userSkillsService.updateUserSkill(skillId, userId, body);
    }
    /**
     * Delete a user skill
     */
    async deleteUserSkill(skillId, userId) {
        await this.userSkillsService.deleteUserSkill(skillId, userId);
        this.setStatus(204);
    }
    /**
     * Get users by skill
     */
    async getUsersBySkill(skillId) {
        return this.userSkillsService.getUsersBySkill(skillId);
    }
    /**
     * Get user skills by proficiency level
     */
    async getUserSkillsByProficiency(userId, proficiencyLevel) {
        return this.userSkillsService.getUserSkillsByProficiency(userId, proficiencyLevel);
    }
};
exports.UserSkillsController = UserSkillsController;
__decorate([
    (0, tsoa_1.Get)('user/{userId}'),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "getUserSkills", null);
__decorate([
    (0, tsoa_1.Get)('{skillId}/user/{userId}'),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, Number]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "getUserSkillById", null);
__decorate([
    (0, tsoa_1.Post)(),
    (0, tsoa_1.SuccessResponse)('201', 'Created'),
    __param(0, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "createUserSkill", null);
__decorate([
    (0, tsoa_1.Put)('{skillId}/user/{userId}'),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Path)()),
    __param(2, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, Number, Object]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "updateUserSkill", null);
__decorate([
    (0, tsoa_1.Delete)('{skillId}/user/{userId}'),
    (0, tsoa_1.SuccessResponse)('204', 'Deleted'),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, Number]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "deleteUserSkill", null);
__decorate([
    (0, tsoa_1.Get)('skill/{skillId}/users'),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "getUsersBySkill", null);
__decorate([
    (0, tsoa_1.Get)('user/{userId}/proficiency/{proficiencyLevel}'),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, String]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "getUserSkillsByProficiency", null);
exports.UserSkillsController = UserSkillsController = __decorate([
    (0, tsoa_1.Route)('user-skills'),
    (0, tsoa_1.Tags)("User Skills"),
    (0, tsoa_1.Security)('BearerAuth')
], UserSkillsController);
