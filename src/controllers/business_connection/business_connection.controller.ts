import {
  Body,
  Controller,
  Post,
  Get,
  Put,
  Path,
  Route,
  SuccessResponse,
  Response,
  Tags,
  Security,
  Query,
} from "tsoa";
import { StatusCodes } from "http-status-codes";
import { ApiError } from "errors/ApiError.js";
import { BusinessConnectionService } from "services/business_connection/business_connection.service.js";
import type {
  InitiateConnectionInput,
  ConnectionIdParams,
  UpdateConnectionStatusInput,
} from "types/buisness_connection/business_connection.schema.js";
import type {
  BusinessConnectionResponse,
  ConnectionSummaryResponse,
} from "types/buisness_connection/business_connection.dto.js";

interface ErrorResponse {
  message: string;
  details?: any;
}

@Route("connections")
@Tags("Business Connections")
@Security("BearerAuth")
export class BusinessConnectionController extends Controller {
  private connectionService = new BusinessConnectionService();

  /**
   * Initiates a new connection request between two businesses.
   * @summary Initiate Connection
   * @returns {BusinessConnectionResponse} The new connection record with status 'pending'.
   */
  @Post()
  @SuccessResponse(StatusCodes.CREATED, "Connection request initiated")
  @Response<ErrorResponse>(StatusCodes.CONFLICT, "Connection already exists")
  @Response<ErrorResponse>(
    StatusCodes.NOT_FOUND,
    "Business or User ID not found",
  )
  public async initiateConnection(
    @Body() body: InitiateConnectionInput,
  ): Promise<BusinessConnectionResponse | ErrorResponse> {
    try {
      const connection = await this.connectionService.initiateConnection(body);
      this.setStatus(StatusCodes.CREATED);
      return connection;
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: "Failed to initiate connection." };
    }
  }

  /**
   * Retrieves all connections associated with a specific business, with optional status filtering.
   * @summary Get Connections for Business
   * @param {number} businessId The ID of the business whose connections to retrieve.
   * @param {string} [status] Optional filter by connection status ('pending', 'active', etc.).
   * @returns {ConnectionSummaryResponse[]} A list of connection summaries.
   */

  @Get("{businessId}")
  public async getConnectionsForBusiness(
    @Path() businessId: number,
    @Query() status?: "pending" | "active" | "rejected" | "inactive",
  ): Promise<ConnectionSummaryResponse[] | ErrorResponse> {
    try {
      if (businessId <= 0) {
        throw new ApiError(StatusCodes.BAD_REQUEST, "Invalid Business ID.");
      }
      return await this.connectionService.getConnectionsForBusiness(
        businessId,
        status,
      );
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: "Failed to retrieve connections." };
    }
  }

  /**
   * Updates the status of an existing connection (e.g., accepting a pending request).
   * @summary Update Connection Status
   * @param {number} connectionId The ID of the connection record to update.
   * @returns {BusinessConnectionResponse} The updated connection record.
   */
  @Put("{connectionId}/status")
  @Response<ErrorResponse>(StatusCodes.NOT_FOUND, "Connection not found")
  @Response<ErrorResponse>(StatusCodes.BAD_REQUEST, "Invalid status update")
  public async updateConnectionStatus(
    @Path() connectionId: ConnectionIdParams["connectionId"],
    @Body() body: UpdateConnectionStatusInput,
  ): Promise<BusinessConnectionResponse | ErrorResponse> {
    try {
      return await this.connectionService.updateConnectionStatus(
        connectionId,
        body,
      );
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: "Failed to update connection status." };
    }
  }

  /**
   * Retrieves a single connection record by ID.
   * @summary Get Connection by ID
   * @param {number} connectionId The ID of the connection record.
   * @returns {BusinessConnectionResponse} The detailed connection record.
   */
  @Get("record/{connectionId}")
  @Response<ErrorResponse>(StatusCodes.NOT_FOUND, "Connection not found")
  public async getConnectionById(
    @Path() connectionId: ConnectionIdParams["connectionId"],
  ): Promise<BusinessConnectionResponse | ErrorResponse> {
    const connection =
      await this.connectionService.getConnectionById(connectionId);
    if (!connection) {
      this.setStatus(StatusCodes.NOT_FOUND);
      return { message: "Business connection not found." };
    }
    return connection;
  }
}
