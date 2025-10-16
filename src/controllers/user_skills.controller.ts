import {
  Controller,
  Get,
  Post,
  Put,
  Delete,
  Path,
  Route,
  Body,
  Security,
  SuccessResponse,
  Tags
} from 'tsoa';
import { UserSkillsService } from '../services/user_skills.service';
import { CreateUserSkillInput, UpdateUserSkillInput } from '../types/user_skills.types';


@Route('user-skills')
@Tags("User Skills")
@Security('BearerAuth')
export class UserSkillsController extends Controller {
  private userSkillsService = new UserSkillsService();

  /**
   * Get all skills for a user
   */
  @Get('user/{userId}')
  public async getUserSkills(@Path() userId: number) {
    return this.userSkillsService.getUserSkills(userId);
  }

  /**
   * Get a specific user skill
   */
  @Get('{skillId}/user/{userId}')
  public async getUserSkillById(
    @Path() skillId: number,
    @Path() userId: number
  ) {
    return this.userSkillsService.getUserSkillById(skillId, userId);
  }

  /**
   * Create a new user skill
   */
  @Post()
  @SuccessResponse('201', 'Created')
  public async createUserSkill(@Body() body: CreateUserSkillInput) {
    const result = await this.userSkillsService.createUserSkill(body);
    this.setStatus(201);
    return result;
  }

  /**
   * Update a user skill
   */
  @Put('{skillId}/user/{userId}')
  public async updateUserSkill(
    @Path() skillId: number,
    @Path() userId: number,
    @Body() body: UpdateUserSkillInput
  ) {
    return this.userSkillsService.updateUserSkill(skillId, userId, body);
  }

  /**
   * Delete a user skill
   */
  @Delete('{skillId}/user/{userId}')
  @SuccessResponse('204', 'Deleted')
  public async deleteUserSkill(
    @Path() skillId: number,
    @Path() userId: number
  ) {
    await this.userSkillsService.deleteUserSkill(skillId, userId);
    this.setStatus(204);
  }

  /**
   * Get users by skill
   */
  @Get('skill/{skillId}/users')
  public async getUsersBySkill(@Path() skillId: number) {
    return this.userSkillsService.getUsersBySkill(skillId);
  }

  /**
   * Get user skills by proficiency level
   */
  @Get('user/{userId}/proficiency/{proficiencyLevel}')
  public async getUserSkillsByProficiency(
    @Path() userId: number,
    @Path() proficiencyLevel: string
  ) {
    return this.userSkillsService.getUserSkillsByProficiency(userId, proficiencyLevel);
  }
}
