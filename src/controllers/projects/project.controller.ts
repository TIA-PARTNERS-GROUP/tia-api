import {
  Body,
  Controller,
  Post,
  Get,
  Put,
  Delete,
  Path,
  Route,
  SuccessResponse,
  Response,
  Tags,
  Security,
} from "tsoa";
import { StatusCodes } from "http-status-codes";
import { ApiError } from "errors/ApiError.js";
import { ProjectService } from "services/projects/project.service.js";
import type {
  CreateProjectInput,
  UpdateProjectInput,
  ProjectIdParams,
  AddMemberInput,
  UpdateMemberRoleInput,
} from "types/projects/project.schema.js";
import type {
  ProjectResponse,
  ProjectMemberResponse,
} from "types/projects/project.dto.js";

interface ErrorResponse {
  message: string;
  details?: any;
}

@Route("projects")
@Tags("Projects")
@Security("BearerAuth")
export class ProjectController extends Controller {
  private projectService = new ProjectService();

  /**
   * Creates a new project.
   * @summary Create Project
   * @returns {ProjectResponse} The newly created project details.
   */
  @Post()
  @SuccessResponse(StatusCodes.CREATED, "Project created successfully")
  @Response<ErrorResponse>(
    StatusCodes.CONFLICT,
    "A project with this name already exists",
  )
  public async createProject(
    @Body() body: CreateProjectInput,
  ): Promise<ProjectResponse | ErrorResponse> {
    try {
      const project = await this.projectService.createProject(body);
      this.setStatus(StatusCodes.CREATED);
      return project;
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: "Failed to create project." };
    }
  }

  /**
   * Retrieve project details by ID.
   * @summary Get Project by ID
   * @param {number} projectId Project's unique ID
   * @returns {ProjectResponse} The detailed project object.
   */
  @Get("{projectId}")
  @Response<ErrorResponse>(StatusCodes.NOT_FOUND, "Project not found")
  public async getProjectById(
    @Path() projectId: ProjectIdParams["projectId"],
  ): Promise<ProjectResponse | ErrorResponse> {
    const project = await this.projectService.getProjectById(projectId);
    if (!project) {
      this.setStatus(StatusCodes.NOT_FOUND);
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
  @Put("{projectId}")
  @Response<ErrorResponse>(StatusCodes.NOT_FOUND, "Project not found")
  @Response<ErrorResponse>(
    StatusCodes.UNPROCESSABLE_ENTITY,
    "Validation failed",
  )
  public async updateProject(
    @Path() projectId: ProjectIdParams["projectId"],
    @Body() body: UpdateProjectInput,
  ): Promise<ProjectResponse | ErrorResponse> {
    try {
      return await this.projectService.updateProject(projectId, body);
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: "Failed to update project." };
    }
  }

  /**
   * Delete a project.
   * @summary Delete Project
   * @param {number} projectId Project's unique ID
   * @returns {void} No content (204) on successful deletion.
   */
  @Delete("{projectId}")
  @SuccessResponse(StatusCodes.NO_CONTENT, "Project deleted successfully")
  @Response<ErrorResponse>(StatusCodes.NOT_FOUND, "Project not found")
  public async deleteProject(
    @Path() projectId: ProjectIdParams["projectId"],
  ): Promise<void | ErrorResponse> {
    const success = await this.projectService.deleteProject(projectId);
    if (!success) {
      this.setStatus(StatusCodes.NOT_FOUND);
      return { message: "Project not found." };
    }
    this.setStatus(StatusCodes.NO_CONTENT);
  }

  /**
   * Retrieves all members associated with a specific project.
   * @summary Get All Project Members
   * @param {number} projectId Project's unique ID
   * @returns {ProjectMemberResponse[]} List of project members.
   * @throws {404} Not Found - Project not found.
   */

  @Get("{projectId}/members")
  @Response<ErrorResponse>(StatusCodes.NOT_FOUND, "Project not found")
  public async getProjectMembers(
    @Path() projectId: ProjectIdParams["projectId"],
  ): Promise<ProjectMemberResponse[] | ErrorResponse> {
    try {
      return await this.projectService.getMembersByProjectId(projectId);
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
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

  @Get("{projectId}/members/{userId}")
  @Response<ErrorResponse>(StatusCodes.NOT_FOUND, "Project member not found")
  public async getProjectMember(
    @Path() projectId: ProjectIdParams["projectId"],
    @Path() userId: number,
  ): Promise<ProjectMemberResponse | ErrorResponse> {
    try {
      const member = await this.projectService.getMemberByKeys(
        projectId,
        userId,
      );
      if (!member) {
        throw new ApiError(StatusCodes.NOT_FOUND, "Project member not found.");
      }
      return member;
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: "Failed to retrieve member." };
    }
  }

  /**
   * Adds a user as a member to the project.
   * @summary Add Project Member
   * @param {number} projectId Project's unique ID
   * @returns {ProjectMemberResponse} The new project member object.
   */
  @Post("{projectId}/members")
  @SuccessResponse(StatusCodes.CREATED, "Member added successfully")
  @Response<ErrorResponse>(StatusCodes.NOT_FOUND, "Project or User not found")
  @Response<ErrorResponse>(StatusCodes.CONFLICT, "User is already a member")
  public async addMember(
    @Path() projectId: ProjectIdParams["projectId"],
    @Body() body: AddMemberInput,
  ): Promise<ProjectMemberResponse | ErrorResponse> {
    try {
      const member = await this.projectService.addMember(projectId, body);
      this.setStatus(StatusCodes.CREATED);
      return member;
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
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
  @Put("{projectId}/members/{userId}")
  @Response<ErrorResponse>(StatusCodes.NOT_FOUND, "Project member not found")
  public async updateMemberRole(
    @Path() projectId: ProjectIdParams["projectId"],
    @Path() userId: number,
    @Body() body: UpdateMemberRoleInput,
  ): Promise<ProjectMemberResponse | ErrorResponse> {
    try {
      return await this.projectService.updateMemberRole(
        projectId,
        userId,
        body,
      );
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
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

  @Delete("{projectId}/members/{userId}")
  @SuccessResponse(StatusCodes.NO_CONTENT, "Member removed successfully")
  @Response<ErrorResponse>(StatusCodes.NOT_FOUND, "Project member not found")
  public async removeMember(
    @Path() projectId: ProjectIdParams["projectId"],
    @Path() userId: number,
  ): Promise<void | ErrorResponse> {
    try {
      const success = await this.projectService.removeMember(projectId, userId);
      if (!success) {
        throw new ApiError(StatusCodes.NOT_FOUND, "Project member not found.");
      }
      this.setStatus(StatusCodes.NO_CONTENT);
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: "Failed to remove member." };
    }
  }
}
