import {
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
  Tags,
  Response,
} from "tsoa";
import {
  CreateBusinessInput,
  UpdateBusinessInput,
} from "types/businesses/businesses.validation.js";
import {
  BusinessResponse,
  BusinessStatsResponse,
  BusinessesFilter,
} from "types/businesses/businesses.dto.js";
import { BusinessesService } from "services/businesses/businesses.service.js";
import { BaseController } from "controllers/base.controller.js";
import { ApiErrorResponse } from "errors/ApiError.js";
import { StatusCodes, ReasonPhrases } from "http-status-codes";

/**
 * Business Management API
 *
 * Provides endpoints for creating, retrieving, updating, and deleting business profiles.
 * Also includes endpoints for managing business status and retrieving statistics.
 *
 * @version 1.0.0
 */
@Route("businesses")
@Tags("Businesses")
@Security("BearerAuth")
@Response<ApiErrorResponse>(
  StatusCodes.INTERNAL_SERVER_ERROR,
  ReasonPhrases.INTERNAL_SERVER_ERROR,
)
export class BusinessesController extends BaseController {
  private readonly businessesService = new BusinessesService();

  /**
   * Get All Businesses with Filtering
   *
   * Retrieves a list of all businesses, with optional query parameters for filtering
   * by type, category, phase, status, and operator user ID.
   *
   * @summary Get a filtered list of businesses
   * @returns {BusinessResponse[]} A list of business profiles.
   */
  @Get()
  public async getBusinesses(
    @Query() business_type?: string,
    @Query() business_category?: string,
    @Query() business_phase?: string,
    @Query() active?: boolean,
    @Query() search?: string,
    @Query() operator_user_id?: number,
  ): Promise<BusinessResponse[] | ApiErrorResponse> {
    try {
      const filters: BusinessesFilter = {
        business_type,
        business_category,
        business_phase,
        active,
        search,
        operator_user_id,
      };
      return await this.businessesService.getBusinesses(filters);
    } catch (error) {
      return this.handleError(error);
    }
  }

  /**
   * Get Business by ID
   *
   * Retrieves detailed information for a specific business by its unique identifier.
   *
   * @summary Get a single business by its ID
   * @param {number} businessId The unique identifier of the business.
   * @returns {BusinessResponse} The requested business profile.
   * @throws {404} Not Found - The business with the specified ID does not exist.
   */
  @Get("{businessId}")
  @Response<ApiErrorResponse>(StatusCodes.NOT_FOUND, "Business not found")
  public async getBusinessById(
    @Path() businessId: number,
  ): Promise<BusinessResponse | ApiErrorResponse> {
    try {
      return await this.businessesService.getBusinessById(businessId);
    } catch (error) {
      return this.handleError(error);
    }
  }

  /**
   * Create a New Business
   *
   * Registers a new business profile in the system. Requires authentication.
   *
   * @summary Create a new business profile
   * @param {CreateBusinessInput} body The data for the new business.
   * @returns {BusinessResponse} The newly created business profile.
   * @throws {401} Unauthorized - Authentication is required.
   * @throws {422} Unprocessable Entity - Validation errors in input data.
   */
  @Post()
  @SuccessResponse(StatusCodes.CREATED, ReasonPhrases.CREATED)
  @Response<ApiErrorResponse>(
    StatusCodes.UNAUTHORIZED,
    ReasonPhrases.UNAUTHORIZED,
  )
  @Response<ApiErrorResponse>(
    StatusCodes.UNPROCESSABLE_ENTITY,
    ReasonPhrases.UNPROCESSABLE_ENTITY,
  )
  public async createBusiness(
    @Body() body: CreateBusinessInput,
  ): Promise<BusinessResponse | ApiErrorResponse> {
    try {
      const result = await this.businessesService.createBusiness(body);
      this.setStatus(StatusCodes.CREATED);
      return result;
    } catch (error) {
      return this.handleError(error);
    }
  }

  /**
   * Update a Business
   *
   * Modifies the details of an existing business profile. Requires authentication.
   *
   * @summary Update an existing business
   * @param {number} businessId The unique identifier of the business to update.
   * @param {UpdateBusinessInput} body The fields to update.
   * @returns {BusinessResponse} The updated business profile.
   * @throws {401} Unauthorized - Authentication is required.
   * @throws {404} Not Found - The business with the specified ID does not exist.
   */
  @Put("{businessId}")
  @Response<ApiErrorResponse>(
    StatusCodes.UNAUTHORIZED,
    ReasonPhrases.UNAUTHORIZED,
  )
  @Response<ApiErrorResponse>(StatusCodes.NOT_FOUND, "Business not found")
  public async updateBusiness(
    @Path() businessId: number,
    @Body() body: UpdateBusinessInput,
  ): Promise<BusinessResponse | ApiErrorResponse> {
    try {
      return await this.businessesService.updateBusiness(businessId, body);
    } catch (error) {
      return this.handleError(error);
    }
  }

  /**
   * Delete a Business
   *
   * Permanently removes a business profile from the system. Requires authentication.
   *
   * @summary Delete a business profile
   * @param {number} businessId The unique identifier of the business to delete.
   * @throws {204} No Content - The business was successfully deleted.
   * @throws {401} Unauthorized - Authentication is required.
   * @throws {404} Not Found - The business with the specified ID does not exist.
   */
  @Delete("{businessId}")
  @SuccessResponse(StatusCodes.NO_CONTENT, "Deleted")
  @Response<ApiErrorResponse>(
    StatusCodes.UNAUTHORIZED,
    ReasonPhrases.UNAUTHORIZED,
  )
  @Response<ApiErrorResponse>(StatusCodes.NOT_FOUND, "Business not found")
  public async deleteBusiness(
    @Path() businessId: number,
  ): Promise<void | ApiErrorResponse> {
    try {
      await this.businessesService.deleteBusiness(businessId);
      this.setStatus(StatusCodes.NO_CONTENT);
    } catch (error) {
      return this.handleError(error);
    }
  }

  /**
   * Get All Businesses for a User
   *
   * Retrieves a list of all businesses associated with a specific user.
   *
   * @summary Get all businesses for a user
   * @param {number} userId The unique identifier of the user.
   * @returns {BusinessResponse[]} A list of business profiles.
   */
  @Get("user/{userId}")
  public async getUserBusinesses(
    @Path() userId: number,
  ): Promise<BusinessResponse[] | ApiErrorResponse> {
    try {
      return await this.businessesService.getUserBusinesses(userId);
    } catch (error) {
      return this.handleError(error);
    }
  }

  /**
   * Toggle Business Active Status
   *
   * Toggles the active status of a business. Requires authentication.
   *
   * @summary Toggle the active status of a business
   * @param {number} businessId The unique identifier of the business.
   * @returns {BusinessResponse} The business profile with the updated status.
   * @throws {401} Unauthorized - Authentication is required.
   * @throws {404} Not Found - The business with the specified ID does not exist.
   */
  @Put("{businessId}/toggle-status")
  @Response<ApiErrorResponse>(
    StatusCodes.UNAUTHORIZED,
    ReasonPhrases.UNAUTHORIZED,
  )
  @Response<ApiErrorResponse>(StatusCodes.NOT_FOUND, "Business not found")
  public async toggleBusinessStatus(
    @Path() businessId: number,
  ): Promise<BusinessResponse | ApiErrorResponse> {
    try {
      return await this.businessesService.toggleBusinessStatus(businessId);
    } catch (error) {
      return this.handleError(error);
    }
  }

  /**
   * Get Business Statistics
   *
   * Retrieves statistics and analytics for a specific business.
   *
   * @summary Get statistics for a business
   * @param {number} businessId The unique identifier of the business.
   * @returns {BusinessStatsResponse} The statistics for the business.
   * @throws {404} Not Found - The business with the specified ID does not exist.
   */
  @Get("{businessId}/stats")
  @Response<ApiErrorResponse>(StatusCodes.NOT_FOUND, "Business not found")
  public async getBusinessStats(
    @Path() businessId: number,
  ): Promise<BusinessStatsResponse | ApiErrorResponse> {
    try {
      return await this.businessesService.getBusinessStats(businessId);
    } catch (error) {
      return this.handleError(error);
    }
  }
}
