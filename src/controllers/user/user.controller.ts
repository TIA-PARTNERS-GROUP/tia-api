import {
  Body,
  Get,
  Path,
  Post,
  Put,
  Delete,
  Route,
  SuccessResponse,
  Response,
  Tags,
  Example,
} from "tsoa";
import type { UserResponse } from "types/user/user.dto.js";
import type {
  UserUpdateInput,
  UserCreateInput,
  UserIdParams,
} from "types/user/user.validation.js";
import { StatusCodes, ReasonPhrases } from "http-status-codes";
import { HttpErrors, ApiErrorResponse } from "errors/ApiError.js";
import * as userService from "services/user/user.services.js";
import { BaseController } from "controllers/base.controller.js";

/**
 * User Management API
 *
 * Provides comprehensive user management capabilities including user registration,
 * profile management, and account administration. All user operations require proper
 * authentication and authorization.
 *
 * @security BearerAuth
 * @version 1.0.0
 */
@Route("users")
@Tags("Users")
@Response<ApiErrorResponse>(StatusCodes.UNAUTHORIZED, "Unauthorized")
@Response<ApiErrorResponse>(StatusCodes.FORBIDDEN, "Forbidden")
@Response<ApiErrorResponse>(
  StatusCodes.INTERNAL_SERVER_ERROR,
  "Internal Server Error",
)
export class UsersController extends BaseController {
  /**
   * Retrieve All Users
   *
   * Returns a paginated list of all registered users in the system.
   * @summary Get all users
   * @description Fetches a complete list of users with their profile information.
   *
   * @returns {UserResponse[]} List of user objects with profile information
   * @throws {401} Unauthorized - Authentication required
   * @throws {403} Forbidden - Insufficient permissions
   * @throws {500} Internal Server Error - Server-side processing error
   */
  @SuccessResponse(StatusCodes.OK, ReasonPhrases.OK)
  @Get("/")
  public async getAllUsers(): Promise<UserResponse[] | ApiErrorResponse> {
    try {
      return await userService.findAllUsers();
    } catch (error) {
      return this.handleError(error);
    }
  }

  /**
   * Get User by ID
   *
   * Retrieves detailed information for a specific user by their unique identifier.
   * Returns complete user profile including contact information and account status.
   *
   * @summary Get user by ID
   * @description Fetches a single user's complete profile information by their unique ID.
   *
   * @param {number} id User's unique identifier (positive integer)
   *
   * @returns {UserResponse} Complete user profile object
   * @throws {401} Unauthorized - Authentication required
   * @throws {404} Not Found - User with specified ID does not exist
   * @throws {500} Internal Server Error - Server-side processing error
   */
  @SuccessResponse(StatusCodes.OK, ReasonPhrases.OK)
  @Response<ApiErrorResponse>(StatusCodes.NOT_FOUND, ReasonPhrases.NOT_FOUND)
  @Get("{id}")
  public async getUserById(
    @Path() id: number,
  ): Promise<UserResponse | ApiErrorResponse> {
    try {
      const user = await userService.findUserById(id);
      if (!user) {
        throw HttpErrors.NotFound(`User with ID ${id} not found.`);
      }
      return user;
    } catch (error) {
      return this.handleError(error);
    }
  }

  /**
   * Create New User
   *
   * Registers a new user in the system. Creates a user account with the provided
   * profile information and securely hashes the password before storage.
   *
   * @summary Create a new user
   * @description Registers a new user account with profile information and credentials.
   * Passwords are automatically hashed using bcrypt before storage.
   *
   * @param {UserCreateInput} requestBody User registration data
   *
   * @returns {UserResponse} Newly created user object (excluding sensitive data)
   * @throws {400} Bad Request - Invalid input data format
   * @throws {409} Conflict - Email address already registered
   * @throws {422} Unprocessable Entity - Validation errors in input data
   * @throws {500} Internal Server Error - Account creation failed
   */
  @SuccessResponse(StatusCodes.CREATED, ReasonPhrases.CREATED)
  @Response<ApiErrorResponse>(
    StatusCodes.UNPROCESSABLE_ENTITY,
    "Validation Failed",
  )
  @Response<ApiErrorResponse>(StatusCodes.CONFLICT, "Email already exists")
  @Example<UserCreateInput>({
    first_name: "Jane",
    last_name: "Doe",
    login_email: "jane.doe@example.com",
    password: "SecurePassword123!",
    contact_email: "jane.contact@example.com",
    contact_phone_no: "0455454321",
    adk_session_id: "123e4567-e89b-12d3-a456-426614174000",
  })
  @Post("/")
  public async createUser(
    @Body() requestBody: UserCreateInput,
  ): Promise<UserResponse | ApiErrorResponse> {
    try {
      const newUser = await userService.createUser(requestBody);
      this.setStatus(StatusCodes.CREATED);
      return newUser;
    } catch (error) {
      return this.handleError(error);
    }
  }

  /**
   * Update User Account
   *
   * Updates an existing user's account information. Supports partial updates where
   * only provided fields are modified. Password updates are automatically hashed.
   *
   * @summary Update user account
   * @description Modifies user account information. All fields are optional - only
   * provided fields will be updated. Password changes are securely handled.
   *
   * @param {number} id User's unique identifier
   * @param {UserUpdateInput} requestBody Fields to update (partial update supported)
   *
   * @returns {UserResponse} Updated user profile object
   * @throws {401} Unauthorized - Authentication required
   * @throws {403} Forbidden - Cannot modify another user's profile
   * @throws {404} Not Found - User with specified ID does not exist
   * @throws {409} Conflict - New email address already in use
   * @throws {422} Unprocessable Entity - Validation errors in input data
   * @throws {500} Internal Server Error - Update operation failed
   */
  @SuccessResponse(StatusCodes.OK, ReasonPhrases.OK)
  @Response<ApiErrorResponse>(StatusCodes.NOT_FOUND, ReasonPhrases.NOT_FOUND)
  @Response<ApiErrorResponse>(
    StatusCodes.UNPROCESSABLE_ENTITY,
    "Validation Failed",
  )
  @Response<ApiErrorResponse>(StatusCodes.CONFLICT, "Email already exists")
  @Response<ApiErrorResponse>(StatusCodes.FORBIDDEN, ReasonPhrases.FORBIDDEN)
  @Put("{id}")
  public async updateUser(
    @Path() id: UserIdParams["id"],
    @Body() requestBody: UserUpdateInput,
  ): Promise<UserResponse | ApiErrorResponse> {
    try {
      return await userService.updateUser(id, requestBody);
    } catch (error) {
      return this.handleError(error);
    }
  }

  /**
   * Delete User Account
   *
   * Permanently removes a user account from the system. This operation is irreversible
   * and will delete all associated user data.
   *
   * @summary Delete user account
   * @description Permanently deletes a user account and associated data. Use with caution
   * as this action cannot be undone.
   *
   * @param {number} id User's unique identifier to delete
   *
   * @throws {204} No Content - User successfully deleted
   * @throws {401} Unauthorized - Authentication required
   * @throws {403} Forbidden - Cannot delete another user's account
   * @throws {404} Not Found - User with specified ID does not exist
   * @throws {500} Internal Server Error - Deletion operation failed
   */
  @SuccessResponse(StatusCodes.NO_CONTENT, ReasonPhrases.NO_CONTENT)
  @Response<ApiErrorResponse>(StatusCodes.NOT_FOUND, ReasonPhrases.NOT_FOUND)
  @Delete("{id}")
  public async deleteUser(
    @Path() id: number,
  ): Promise<void | ApiErrorResponse> {
    try {
      await userService.deleteUser(id);
      this.setStatus(StatusCodes.NO_CONTENT);
    } catch (error) {
      return this.handleError(error);
    }
  }
}
