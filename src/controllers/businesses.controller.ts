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
  SuccessResponse
} from 'tsoa';
import { BusinessesService } from '../services/businesses.service';
import { CreateBusinessInput, UpdateBusinessInput } from '../types/businesses.types';


@Route('businesses')
export class BusinessesController extends Controller {
  private businessesService = new BusinessesService();

  @Get()
  public async getBusinesses(
    @Query() business_type?: string,
    @Query() business_category?: string,
    @Query() business_phase?: string,
    @Query() active?: boolean,
    @Query() search?: string,
    @Query() operator_user_id?: number
  ) {
    const filters: any = {};
    if (business_type) filters.business_type = business_type;
    if (business_category) filters.business_category = business_category;
    if (business_phase) filters.business_phase = business_phase;
    if (active !== undefined) filters.active = active;
    if (search) filters.search = search;
    if (operator_user_id) filters.operator_user_id = operator_user_id;

    return this.businessesService.getBusinesses(filters);
  }

  @Get('{businessId}')
  public async getBusinessById(@Path() businessId: number) {
    return this.businessesService.getBusinessById(businessId);
  }

  @Post()
  @Security("BearerAuth")
  @SuccessResponse('201', 'Created')
  public async createBusiness(@Body() body: CreateBusinessInput) {
    const result = await this.businessesService.createBusiness(body);
    this.setStatus(201);
    return result;
  }

  @Put('{businessId}')
  @Security("BearerAuth")
  public async updateBusiness(
    @Path() businessId: number,
    @Body() body: UpdateBusinessInput
  ) {
    return this.businessesService.updateBusiness(businessId, body);
  }

  @Delete('{businessId}')
  @Security("BearerAuth")
  @SuccessResponse('204', 'Deleted')
  public async deleteBusiness(@Path() businessId: number) {
    await this.businessesService.deleteBusiness(businessId);
    this.setStatus(204);
  }

  @Get('user/{userId}')
  public async getUserBusinesses(@Path() userId: number) {
    return this.businessesService.getUserBusinesses(userId);
  }

  @Put('{businessId}/toggle-status')
  @Security("BearerAuth")
  public async toggleBusinessStatus(@Path() businessId: number) {
    return this.businessesService.toggleBusinessStatus(businessId);
  }

  @Get('{businessId}/stats')
  public async getBusinessStats(@Path() businessId: number) {
    return this.businessesService.getBusinessStats(businessId);
  }
}
