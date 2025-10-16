import {
  Body,
  Controller,
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
} from 'tsoa';
import { StatusCodes, ReasonPhrases } from 'http-status-codes';
import { ApiError } from '../errors/ApiError';
import * as userService from '@services/user.services';
import type {
  UserResponse,
  UserCreateRequest,
  UserUpdateRequest
} from '../types/user.dto.js';
import type { UserUpdateParams } from '../types/user.validation';

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
export class UsersController extends Controller {
  /**
   * Retrieve All Users
   * 
   * Returns a paginated list of all registered users in the system.
   * @summary Get all users
   * @description Fetches a complete list of users with their profile information. 
   * 
   * @isInt @minimum 1 @default 1 page Page number for pagination
   * @isInt @minimum 1 @maximum 100 @default 20 limit Number of users per page
   * 
   * @example page 1
   * @example limit 20
   * 
   * @returns {UserResponse[]} List of user objects with profile information
   * @throws {401} Unauthorized - Authentication required
   * @throws {403} Forbidden - Insufficient permissions
   * @throws {500} Internal Server Error - Server-side processing error
   */
  @SuccessResponse(StatusCodes.OK, ReasonPhrases.OK)
  @Get("/")
  public async getAllUsers(): Promise<UserResponse[]> {
    return await userService.findAllUsers();
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
   * @example id 123
   * 
   * @returns {UserResponse} Complete user profile object
   * @throws {401} Unauthorized - Authentication required
   * @throws {404} Not Found - User with specified ID does not exist
   * @throws {500} Internal Server Error - Server-side processing error
   */
  @SuccessResponse(StatusCodes.OK, ReasonPhrases.OK)
  @Response(StatusCodes.NOT_FOUND, ReasonPhrases.NOT_FOUND)
  @Get("{id}")
  public async getUserById(@Path() id: number): Promise<UserResponse | null> {
    const user = await userService.findUserById(id);
    if (!user) {
      this.setStatus(StatusCodes.NOT_FOUND);
    }
    return user;
  }

  /**
   * Create New User
   * 
   * Registers a new user in the system. Creates a user account with the provided
   * profile information and securely hashes the password before storage.
   * 
   * @summary Create a new user
   * @description Registers a new user account with profile information and credentials.
   *              Passwords are automatically hashed using bcrypt before storage.
   * 
   * @param {UserCreateRequest} requestBody User registration data
   * 
   * @example requestBody {
   *   "first_name": "Jane",
   *   "last_name": "Doe",
   *   "login_email": "jane.doe@example.com",
   *   "password": "SecurePassword123!",
   *   "contact_email": "jane.contact@example.com",
   *   "contact_phone_no": "0455454321",
   *   "adk_session_id": "123e4567-e89b-12d3-a456-426614174000"
   * }
   * 
   * @returns {UserResponse} Newly created user object (excluding sensitive data)
   * @throws {400} Bad Request - Invalid input data format
   * @throws {409} Conflict - Email address already registered
   * @throws {422} Unprocessable Entity - Validation errors in input data
   * @throws {500} Internal Server Error - Account creation failed
   */
  @SuccessResponse(StatusCodes.CREATED, ReasonPhrases.CREATED)
  @Response(StatusCodes.UNPROCESSABLE_ENTITY, "Validation Failed")
  @Response(StatusCodes.CONFLICT, "Email already exists")
  @Example<UserCreateRequest>({
    first_name: "Jane",
    last_name: "Doe",
    login_email: "jane.doe@example.com",
    password: "SecurePassword123!",
    contact_email: "jane.contact@example.com",
    contact_phone_no: "0455454321",
    adk_session_id: "123e4567-e89b-12d3-a456-426614174000"
  })
  @Post("/")
  public async createUser(@Body() requestBody: UserCreateRequest): Promise<UserResponse | { message: string; details?: any }> {
    try {
      const result = await userService.createUser(requestBody);
      this.setStatus(StatusCodes.CREATED);
      return result;
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message, details: error.details };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: 'An unexpected error occurred', details: null };
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
   *              provided fields will be updated. Password changes are securely handled.
   * 
   * @param {number} id User's unique identifier
   * @param {UserUpdateRequest} requestBody Fields to update (partial update supported)
   * 
   * @example id 123
   * @example requestBody {
   *   "first_name": "Jane",
   *   "last_name": "Smith",
   *   "login_email": "jane.smith@example.com",
   *   "password": "NewSecurePassword456!",
   *   "contact_email": "jane.smith@example.com",
   *   "contact_phone_no": "0455454321",
   *   "email_verified": true,
   *   "active": true
   * }
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
  @Response(StatusCodes.NOT_FOUND, ReasonPhrases.NOT_FOUND)
  @Response(StatusCodes.UNPROCESSABLE_ENTITY, "Validation Failed")
  @Response(StatusCodes.CONFLICT, "Email already exists")
  @Example<UserUpdateRequest>({
    first_name: "Jane",
    last_name: "Smith",
    login_email: "jane.smith@example.com",
    password: "NewSecurePassword456!",
    contact_email: "jane.smith@example.com",
    contact_phone_no: "0455454321",
    email_verified: true,
    active: true
  })
  @Put("{id}")
  public async updateUser(@Path() id: number, @Body() requestBody: UserUpdateRequest): Promise<UserResponse | { message: string; details?: any }> {
    try {
      const serviceParams: UserUpdateParams = {
        first_name: requestBody.first_name,
        last_name: requestBody.last_name,
        login_email: requestBody.login_email,
        password: requestBody.password,
        contact_email: requestBody.contact_email ?? undefined,
        contact_phone_no: requestBody.contact_phone_no ?? undefined,
        adk_session_id: requestBody.adk_session_id ?? undefined,
        email_verified: requestBody.email_verified,
        active: requestBody.active,
      };

      Object.keys(serviceParams).forEach(key => {
        if (serviceParams[key as keyof UserUpdateParams] === undefined) {
          delete serviceParams[key as keyof UserUpdateParams];
        }
      });

      const updatedUser = await userService.updateUser(id, serviceParams);
      if (!updatedUser) {
        this.setStatus(StatusCodes.NOT_FOUND);
        return { message: 'User not found', details: null };
      }
      return updatedUser;
    } catch (error) {
      if (error instanceof ApiError) {
        this.setStatus(error.statusCode);
        return { message: error.message, details: error.details };
      }
      this.setStatus(StatusCodes.INTERNAL_SERVER_ERROR);
      return { message: 'An unexpected error occurred', details: null };
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
   *              as this action cannot be undone.
   * 
   * @param {number} id User's unique identifier to delete
   * 
   * @example id 123
   * 
   * @throws {204} No Content - User successfully deleted
   * @throws {401} Unauthorized - Authentication required
   * @throws {403} Forbidden - Cannot delete another user's account
   * @throws {404} Not Found - User with specified ID does not exist
   * @throws {500} Internal Server Error - Deletion operation failed
   */
  @SuccessResponse(StatusCodes.NO_CONTENT, ReasonPhrases.NO_CONTENT)
  @Response(StatusCodes.NOT_FOUND, ReasonPhrases.NOT_FOUND)
  @Delete("{id}")
  public async deleteUser(@Path() id: number): Promise<void> {
    const deletedUser = await userService.deleteUser(id);
    if (!deletedUser) {
      this.setStatus(StatusCodes.NOT_FOUND);
      return;
    }
    this.setStatus(StatusCodes.NO_CONTENT);
  }
}
