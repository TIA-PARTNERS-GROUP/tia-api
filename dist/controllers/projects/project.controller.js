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
exports.ProjectController = void 0;
const tsoa_1 = require("tsoa");
const http_status_codes_1 = require("http-status-codes");
const ApiError_1 = require("../../errors/ApiError");
const project_service_1 = require("../../services/projects/project.service");
let ProjectController = class ProjectController extends tsoa_1.Controller {
    projectService = new project_service_1.ProjectService();
    /**
     * Creates a new project.
     * @summary Create Project
     * @returns {ProjectResponse} The newly created project details.
     */
    async createProject(body) {
        try {
            const project = await this.projectService.createProject(body);
            this.setStatus(http_status_codes_1.StatusCodes.CREATED);
            return project;
        }
        catch (error) {
            if (error instanceof ApiError_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
            return { message: "Failed to create project." };
        }
    }
    /**
     * Retrieve project details by ID.
     * @summary Get Project by ID
     * @param {number} projectId Project's unique ID
     * @returns {ProjectResponse} The detailed project object.
     */
    async getProjectById(projectId) {
        const project = await this.projectService.getProjectById(projectId);
        if (!project) {
            this.setStatus(http_status_codes_1.StatusCodes.NOT_FOUND);
            return { message: "Project not found." };
        }
        return project;
    }
    /**
     * Update project details (partial updates allowed).
     * @summary Update Project
     * @param {number} projectId Project's unique ID
     * @returns {ProjectResponse} The updated project details.
     */
    async updateProject(projectId, body) {
        try {
            return await this.projectService.updateProject(projectId, body);
        }
        catch (error) {
            if (error instanceof ApiError_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
            return { message: "Failed to update project." };
        }
    }
    /**
     * Delete a project.
     * @summary Delete Project
     * @param {number} projectId Project's unique ID
     * @returns {void} No content (204) on successful deletion.
     */
    async deleteProject(projectId) {
        const success = await this.projectService.deleteProject(projectId);
        if (!success) {
            this.setStatus(http_status_codes_1.StatusCodes.NOT_FOUND);
            return { message: "Project not found." };
        }
        this.setStatus(http_status_codes_1.StatusCodes.NO_CONTENT);
    }
    /**
     * Retrieves all members associated with a specific project.
     * @summary Get All Project Members
     * @param {number} projectId Project's unique ID
     * @returns {ProjectMemberResponse[]} List of project members.
     * @throws {404} Not Found - Project not found.
     */
    async getProjectMembers(projectId) {
        try {
            return await this.projectService.getMembersByProjectId(projectId);
        }
        catch (error) {
            if (error instanceof ApiError_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
            return { message: "Failed to retrieve members." };
        }
    }
    /**
     * Retrieves a specific member's details within a project.
     * @summary Get Project Member by ID
     * @param {number} projectId Project's unique ID
     * @param {number} userId User's unique ID
     * @returns {ProjectMemberResponse} The specific project member object.
     * @throws {404} Not Found - Project or member not found.
     */
    async getProjectMember(projectId, userId) {
        try {
            const member = await this.projectService.getMemberByKeys(projectId, userId);
            if (!member) {
                throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, "Project member not found.");
            }
            return member;
        }
        catch (error) {
            if (error instanceof ApiError_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
            return { message: "Failed to retrieve member." };
        }
    }
    /**
     * Adds a user as a member to the project.
     * @summary Add Project Member
     * @param {number} projectId Project's unique ID
     * @returns {ProjectMemberResponse} The new project member object.
     */
    async addMember(projectId, body) {
        try {
            const member = await this.projectService.addMember(projectId, body);
            this.setStatus(http_status_codes_1.StatusCodes.CREATED);
            return member;
        }
        catch (error) {
            if (error instanceof ApiError_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
            return { message: "Failed to add member." };
        }
    }
    /**
     * Updates the role of an existing project member.
     * @summary Update Project Member Role
     * @param {number} projectId Project's unique ID
     * @param {number} userId User's unique ID
     * @returns {ProjectMemberResponse} The updated project member object.
     */
    async updateMemberRole(projectId, userId, body) {
        try {
            return await this.projectService.updateMemberRole(projectId, userId, body);
        }
        catch (error) {
            if (error instanceof ApiError_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
            return { message: "Failed to update member role." };
        }
    }
    /**
     * Removes a member from the project.
     * @summary Remove Project Member
     * @param {number} projectId Project's unique ID
     * @param {number} userId User's unique ID to remove
     * @returns {void} No content (204) on successful removal.
     * @throws {404} Not Found - Project or member not found.
     */
    async removeMember(projectId, userId) {
        try {
            const success = await this.projectService.removeMember(projectId, userId);
            if (!success) {
                throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, "Project member not found.");
            }
            this.setStatus(http_status_codes_1.StatusCodes.NO_CONTENT);
        }
        catch (error) {
            if (error instanceof ApiError_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
            return { message: "Failed to remove member." };
        }
    }
};
exports.ProjectController = ProjectController;
__decorate([
    (0, tsoa_1.Post)(),
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.CREATED, "Project created successfully"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.CONFLICT, "A project with this name already exists"),
    __param(0, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], ProjectController.prototype, "createProject", null);
__decorate([
    (0, tsoa_1.Get)("{projectId}"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Project not found"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], ProjectController.prototype, "getProjectById", null);
__decorate([
    (0, tsoa_1.Put)("{projectId}"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Project not found"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNPROCESSABLE_ENTITY, "Validation failed"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object, Object]),
    __metadata("design:returntype", Promise)
], ProjectController.prototype, "updateProject", null);
__decorate([
    (0, tsoa_1.Delete)("{projectId}"),
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.NO_CONTENT, "Project deleted successfully"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Project not found"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], ProjectController.prototype, "deleteProject", null);
__decorate([
    (0, tsoa_1.Get)("{projectId}/members"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Project not found"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], ProjectController.prototype, "getProjectMembers", null);
__decorate([
    (0, tsoa_1.Get)("{projectId}/members/{userId}"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Project member not found"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object, Number]),
    __metadata("design:returntype", Promise)
], ProjectController.prototype, "getProjectMember", null);
__decorate([
    (0, tsoa_1.Post)("{projectId}/members"),
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.CREATED, "Member added successfully"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Project or User not found"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.CONFLICT, "User is already a member"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object, Object]),
    __metadata("design:returntype", Promise)
], ProjectController.prototype, "addMember", null);
__decorate([
    (0, tsoa_1.Put)("{projectId}/members/{userId}"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Project member not found"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Path)()),
    __param(2, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object, Number, Object]),
    __metadata("design:returntype", Promise)
], ProjectController.prototype, "updateMemberRole", null);
__decorate([
    (0, tsoa_1.Delete)("{projectId}/members/{userId}"),
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.NO_CONTENT, "Member removed successfully"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, "Project member not found"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object, Number]),
    __metadata("design:returntype", Promise)
], ProjectController.prototype, "removeMember", null);
exports.ProjectController = ProjectController = __decorate([
    (0, tsoa_1.Route)("projects"),
    (0, tsoa_1.Tags)("Projects"),
    (0, tsoa_1.Security)("BearerAuth")
], ProjectController);
