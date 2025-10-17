import {
  Get,
  Post,
  Put,
  Delete,
  Path,
  Route,
  Body,
  Security,
  SuccessResponse,
  Tags,
  Response,
} from "tsoa";
import { UserSkillsService } from "services/user_skills/user_skills.service.js";
import {
  CreateUserSkillInput,
  UpdateUserSkillInput,
  UserSkillResponse,
} from "types/user_skills/user_skills.dto.js";
import { BaseController } from "controllers/base.controller.js";
import { ApiErrorResponse } from "errors/ApiError.js";
import { StatusCodes, ReasonPhrases } from "http-status-codes";
import { UserResponse } from "types/user/user.dto.js";

/**
 * User Skills Management API
 *
 * Manages the relationships between users and skills, including proficiency levels.
 * Allows for tracking and updating the skills repertoire of each user.
 *
 * @security BearerAuth
 * @version 1.0.0
 */
@Route("user-skills")
@Tags("User Skills")
@Security("BearerAuth")
@Response<ApiErrorResponse>(
  StatusCodes.UNAUTHORIZED,
  ReasonPhrases.UNAUTHORIZED,
)
@Response<ApiErrorResponse>(
  StatusCodes.INTERNAL_SERVER_ERROR,
  ReasonPhrases.INTERNAL_SERVER_ERROR,
)
export class UserSkillsController extends BaseController {
  private readonly userSkillsService = new UserSkillsService();

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
  @Get("user/{userId}")
  public async getUserSkills(
    @Path() userId: number,
  ): Promise<UserSkillResponse[] | ApiErrorResponse> {
    try {
      return await this.userSkillsService.getUserSkills(userId);
    } catch (error) {
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
  @Get("{skillId}/user/{userId}")
  @Response<ApiErrorResponse>(
    StatusCodes.NOT_FOUND,
    "User-skill relationship not found",
  )
  public async getUserSkillById(
    @Path() skillId: number,
    @Path() userId: number,
  ): Promise<UserSkillResponse | ApiErrorResponse> {
    try {
      return await this.userSkillsService.getUserSkillById(skillId, userId);
    } catch (error) {
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
  @Post()
  @SuccessResponse(StatusCodes.CREATED, ReasonPhrases.CREATED)
  @Response<ApiErrorResponse>(
    StatusCodes.CONFLICT,
    "User already has this skill",
  )
  @Response<ApiErrorResponse>(
    StatusCodes.UNPROCESSABLE_ENTITY,
    ReasonPhrases.UNPROCESSABLE_ENTITY,
  )
  public async createUserSkill(
    @Body() body: CreateUserSkillInput,
  ): Promise<UserSkillResponse | ApiErrorResponse> {
    try {
      const result = await this.userSkillsService.createUserSkill(body);
      this.setStatus(StatusCodes.CREATED);
      return result;
    } catch (error) {
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
  @Put("{skillId}/user/{userId}")
  @Response<ApiErrorResponse>(
    StatusCodes.NOT_FOUND,
    "User-skill relationship not found",
  )
  public async updateUserSkill(
    @Path() skillId: number,
    @Path() userId: number,
    @Body() body: UpdateUserSkillInput,
  ): Promise<UserSkillResponse | ApiErrorResponse> {
    try {
      return await this.userSkillsService.updateUserSkill(
        skillId,
        userId,
        body,
      );
    } catch (error) {
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
  @Delete("{skillId}/user/{userId}")
  @SuccessResponse(StatusCodes.NO_CONTENT, "Deleted")
  @Response<ApiErrorResponse>(
    StatusCodes.NOT_FOUND,
    "User-skill relationship not found",
  )
  public async deleteUserSkill(
    @Path() skillId: number,
    @Path() userId: number,
  ): Promise<void | ApiErrorResponse> {
    try {
      await this.userSkillsService.deleteUserSkill(skillId, userId);
      this.setStatus(StatusCodes.NO_CONTENT);
    } catch (error) {
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
  @Get("skill/{skillId}/users")
  public async getUsersBySkill(
    @Path() skillId: number,
  ): Promise<UserResponse[] | ApiErrorResponse> {
    try {
      return await this.userSkillsService.getUsersBySkill(skillId);
    } catch (error) {
      return this.handleError(error);
    }
  }
}
