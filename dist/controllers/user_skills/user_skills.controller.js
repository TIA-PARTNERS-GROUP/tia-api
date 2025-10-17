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
const user_skills_service_1 = require("../../services/user_skills/user_skills.service");
const base_controller_1 = require("../../controllers/base.controller");
const http_status_codes_1 = require("http-status-codes");
/**
 * User Skills Management API
 *
 * Manages the relationships between users and skills, including proficiency levels.
 * Allows for tracking and updating the skills repertoire of each user.
 *
 * @security BearerAuth
 * @version 1.0.0
 */
let UserSkillsController = class UserSkillsController extends base_controller_1.BaseController {
    userSkillsService = new user_skills_service_1.UserSkillsService();
    /**
     * Get All Skills for a User
     *
     * Retrieves a list of all skills associated with a specific user, including their
     * proficiency level and the date it was last updated.
     *
     * @summary Get all skills for a specific user
     * @param {number} userId The unique identifier of the user.
     * @returns {UserSkillResponse[]} A list of skills associated with the user.
     */
    async getUserSkills(userId) {
        try {
            return await this.userSkillsService.getUserSkills(userId);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Get a Specific User-Skill Relationship
     *
     * Retrieves a single skill entry for a specific user by their respective IDs.
     *
     * @summary Get a specific user-skill entry
     * @param {number} skillId The unique identifier for the skill.
     * @param {number} userId The unique identifier for the user.
     * @returns {UserSkillResponse} The specific user-skill relationship object.
     * @throws {404} Not Found - The specified user or skill association does not exist.
     */
    async getUserSkillById(skillId, userId) {
        try {
            return await this.userSkillsService.getUserSkillById(skillId, userId);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Add a Skill to a User
     *
     * Creates a new association between a user and a skill, including the proficiency level.
     *
     * @summary Add a new skill to a user's profile
     * @param {CreateUserSkillInput} body The user ID, skill ID, and proficiency level.
     * @returns {UserSkillResponse} The newly created user-skill relationship object.
     * @throws {409} Conflict - The user is already associated with this skill.
     * @throws {422} Unprocessable Entity - Invalid input data.
     */
    async createUserSkill(body) {
        try {
            const result = await this.userSkillsService.createUserSkill(body);
            this.setStatus(http_status_codes_1.StatusCodes.CREATED);
            return result;
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Update a User's Skill Proficiency
     *
     * Modifies the proficiency level or other details of an existing user-skill association.
     *
     * @summary Update a user's skill details
     * @param {number} skillId The unique identifier for the skill.
     * @param {number} userId The unique identifier for the user.
     * @param {UpdateUserSkillInput} body The fields to update, such as proficiency level.
     * @returns {UserSkillResponse} The updated user-skill relationship object.
     * @throws {404} Not Found - The specified user or skill association does not exist.
     */
    async updateUserSkill(skillId, userId, body) {
        try {
            return await this.userSkillsService.updateUserSkill(skillId, userId, body);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Remove a Skill from a User
     *
     * Permanently deletes the association between a user and a skill.
     *
     * @summary Remove a skill from a user's profile
     * @param {number} skillId The unique identifier for the skill to remove.
     * @param {number} userId The unique identifier of the user.
     * @throws {204} No Content - The skill was successfully removed.
     * @throws {404} Not Found - The specified user or skill association does not exist.
     */
    async deleteUserSkill(skillId, userId) {
        try {
            await this.userSkillsService.deleteUserSkill(skillId, userId);
            this.setStatus(http_status_codes_1.StatusCodes.NO_CONTENT);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Find Users by Skill
     *
     * Retrieves a list of all users who possess a specific skill.
     *
     * @summary Get all users who have a specific skill
     * @param {number} skillId The unique identifier of the skill.
     * @returns {UserResponse[]} A list of user profiles.
     */
    async getUsersBySkill(skillId) {
        try {
            return await this.userSkillsService.getUsersBySkill(skillId);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
    /**
     * Get User Skills by Proficiency
     *
     * Retrieves a list of a user's skills that match a specific proficiency level.
     *
     * @summary Filter a user's skills by proficiency
     * @param {number} userId The unique identifier of the user.
     * @param {string} proficiencyLevel The proficiency level to filter by (e.g., "Beginner", "Expert").
     * @returns {UserSkillResponse[]} A list of skills matching the proficiency level.
     */
    async getUserSkillsByProficiency(userId, proficiencyLevel) {
        try {
            return this.userSkillsService.getUserSkillsByProficiency(userId, proficiencyLevel);
        }
        catch (error) {
            return this.handleError(error);
        }
    }
};
exports.UserSkillsController = UserSkillsController;
__decorate([
    (0, tsoa_1.Get)("user/{userId}"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "getUserSkills", null);
__decorate([
    (0, tsoa_1.Get)("{skillId}/user/{userId}"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "User-skill relationship not found"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, Number]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "getUserSkillById", null);
__decorate([
    (0, tsoa_1.Post)(),
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.CREATED, http_status_codes_1.ReasonPhrases.CREATED),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.CONFLICT, "User already has this skill"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNPROCESSABLE_ENTITY, http_status_codes_1.ReasonPhrases.UNPROCESSABLE_ENTITY),
    __param(0, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "createUserSkill", null);
__decorate([
    (0, tsoa_1.Put)("{skillId}/user/{userId}"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "User-skill relationship not found"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Path)()),
    __param(2, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, Number, Object]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "updateUserSkill", null);
__decorate([
    (0, tsoa_1.Delete)("{skillId}/user/{userId}"),
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.NO_CONTENT, "Deleted"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "User-skill relationship not found"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, Number]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "deleteUserSkill", null);
__decorate([
    (0, tsoa_1.Get)("skill/{skillId}/users"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "getUsersBySkill", null);
__decorate([
    (0, tsoa_1.Get)("user/{userId}/proficiency/{proficiencyLevel}"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, String]),
    __metadata("design:returntype", Promise)
], UserSkillsController.prototype, "getUserSkillsByProficiency", null);
exports.UserSkillsController = UserSkillsController = __decorate([
    (0, tsoa_1.Route)("user-skills"),
    (0, tsoa_1.Tags)("User Skills"),
    (0, tsoa_1.Security)("BearerAuth"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNAUTHORIZED, http_status_codes_1.ReasonPhrases.UNAUTHORIZED),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR, http_status_codes_1.ReasonPhrases.INTERNAL_SERVER_ERROR)
], UserSkillsController);
