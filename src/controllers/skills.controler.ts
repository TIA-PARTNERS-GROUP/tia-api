import {
  Controller,
  Get,
  Post,
  Put,
  Delete,
  Path,
  Route,
  Body,
  Query,
  Security,
  SuccessResponse,
  Tags
} from 'tsoa';
import { SkillsService } from '../services/skills.service';
import { CreateSkillInput, UpdateSkillInput, SkillsFilter } from '../types/skills.types';

@Route('skills')
@Tags("Skills")
export class SkillsController extends Controller {
  private skillsService = new SkillsService();

  /**
   * Get all skills with optional filtering
   */
  @Get()
  public async getSkills(
    @Query() category?: string,
    @Query() active?: boolean,
    @Query() search?: string
  ) {
    const filters: SkillsFilter = {};
    if (category) filters.category = category;
    if (active !== undefined) filters.active = active;
    if (search) filters.search = search;

    return this.skillsService.getSkills(filters);
  }

  /**
   * Get a specific skill by ID
   */
  @Get('{skillId}')
  public async getSkillById(@Path() skillId: number) {
    return this.skillsService.getSkillById(skillId);
  }

  /**
   * Get a skill by name
   */
  @Get('name/{name}')
  public async getSkillByName(@Path() name: string) {
    return this.skillsService.getSkillByName(name);
  }

  /**
   * Create a new skill
   */
  @Post()
  @Security('BearerAuth')
  @SuccessResponse('201', 'Created')
  public async createSkill(@Body() body: CreateSkillInput) {
    const result = await this.skillsService.createSkill(body);
    this.setStatus(201);
    return result;
  }

  /**
   * Update a skill
   */
  @Put('{skillId}')
  @Security('BearerAuth')
  public async updateSkill(
    @Path() skillId: number,
    @Body() body: UpdateSkillInput
  ) {
    return this.skillsService.updateSkill(skillId, body);
  }

  /**
   * Delete a skill
   */
  @Delete('{skillId}')
  @Security('BearerAuth')
  @SuccessResponse('204', 'Deleted')
  public async deleteSkill(@Path() skillId: number) {
    await this.skillsService.deleteSkill(skillId);
    this.setStatus(204);
  }

  /**
   * Get skills by category
   */
  @Get('category/{category}')
  public async getSkillsByCategory(@Path() category: string) {
    return this.skillsService.getSkillsByCategory(category);
  }

  /**
   * Get all unique categories
   */
  @Get('categories/all')
  public async getSkillCategories() {
    return this.skillsService.getSkillCategories();
  }

  /**
   * Toggle skill active status
   */
  @Put('{skillId}/toggle-status')
  @Security('BearerAuth')
  public async toggleSkillStatus(@Path() skillId: number) {
    return this.skillsService.toggleSkillStatus(skillId);
  }

  /**
   * Get popular skills
   */
  @Get('popular')
  public async getPopularSkills(@Query() limit?: number) {
    return this.skillsService.getPopularSkills(limit);
  }

  /**
   * Search skills
   */
  @Get('search/{query}')
  public async searchSkills(
    @Path() query: string,
    @Query() limit?: number
  ) {
    return this.skillsService.searchSkills(query, limit);
  }
}
