"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var __param = (this && this.__param) || function (paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.UsersController = void 0;
const tsoa_1 = require("tsoa");
const http_status_codes_1 = require("http-status-codes");
const ApiError_1 = require("../errors/ApiError");
const userService = __importStar(require("../services/user.services"));
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
let UsersController = class UsersController extends tsoa_1.Controller {
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
    async getAllUsers() {
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
    async getUserById(id) {
        const user = await userService.findUserById(id);
        if (!user) {
            this.setStatus(http_status_codes_1.StatusCodes.NOT_FOUND);
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
    async createUser(requestBody) {
        try {
            const result = await userService.createUser(requestBody);
            this.setStatus(http_status_codes_1.StatusCodes.CREATED);
            return result;
        }
        catch (error) {
            if (error instanceof ApiError_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message, details: error.details };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
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
    async updateUser(id, requestBody) {
        try {
            const serviceParams = {
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
                if (serviceParams[key] === undefined) {
                    delete serviceParams[key];
                }
            });
            const updatedUser = await userService.updateUser(id, serviceParams);
            if (!updatedUser) {
                this.setStatus(http_status_codes_1.StatusCodes.NOT_FOUND);
                return { message: 'User not found', details: null };
            }
            return updatedUser;
        }
        catch (error) {
            if (error instanceof ApiError_1.ApiError) {
                this.setStatus(error.statusCode);
                return { message: error.message, details: error.details };
            }
            this.setStatus(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR);
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
    async deleteUser(id) {
        const deletedUser = await userService.deleteUser(id);
        if (!deletedUser) {
            this.setStatus(http_status_codes_1.StatusCodes.NOT_FOUND);
            return;
        }
        this.setStatus(http_status_codes_1.StatusCodes.NO_CONTENT);
    }
};
exports.UsersController = UsersController;
__decorate([
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.OK, http_status_codes_1.ReasonPhrases.OK),
    (0, tsoa_1.Get)("/"),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", []),
    __metadata("design:returntype", Promise)
], UsersController.prototype, "getAllUsers", null);
__decorate([
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.OK, http_status_codes_1.ReasonPhrases.OK),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, http_status_codes_1.ReasonPhrases.NOT_FOUND),
    (0, tsoa_1.Get)("{id}"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], UsersController.prototype, "getUserById", null);
__decorate([
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.CREATED, http_status_codes_1.ReasonPhrases.CREATED),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNPROCESSABLE_ENTITY, "Validation Failed"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.CONFLICT, "Email already exists"),
    (0, tsoa_1.Example)({
        first_name: "Jane",
        last_name: "Doe",
        login_email: "jane.doe@example.com",
        password: "SecurePassword123!",
        contact_email: "jane.contact@example.com",
        contact_phone_no: "0455454321",
        adk_session_id: "123e4567-e89b-12d3-a456-426614174000"
    }),
    (0, tsoa_1.Post)("/"),
    __param(0, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Object]),
    __metadata("design:returntype", Promise)
], UsersController.prototype, "createUser", null);
__decorate([
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.OK, http_status_codes_1.ReasonPhrases.OK),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, http_status_codes_1.ReasonPhrases.NOT_FOUND),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.UNPROCESSABLE_ENTITY, "Validation Failed"),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.CONFLICT, "Email already exists"),
    (0, tsoa_1.Example)({
        first_name: "Jane",
        last_name: "Smith",
        login_email: "jane.smith@example.com",
        password: "NewSecurePassword456!",
        contact_email: "jane.smith@example.com",
        contact_phone_no: "0455454321",
        email_verified: true,
        active: true
    }),
    (0, tsoa_1.Put)("{id}"),
    __param(0, (0, tsoa_1.Path)()),
    __param(1, (0, tsoa_1.Body)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, Object]),
    __metadata("design:returntype", Promise)
], UsersController.prototype, "updateUser", null);
__decorate([
    (0, tsoa_1.SuccessResponse)(http_status_codes_1.StatusCodes.NO_CONTENT, http_status_codes_1.ReasonPhrases.NO_CONTENT),
    (0, tsoa_1.Response)(http_status_codes_1.StatusCodes.NOT_FOUND, http_status_codes_1.ReasonPhrases.NOT_FOUND),
    (0, tsoa_1.Delete)("{id}"),
    __param(0, (0, tsoa_1.Path)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number]),
    __metadata("design:returntype", Promise)
], UsersController.prototype, "deleteUser", null);
exports.UsersController = UsersController = __decorate([
    (0, tsoa_1.Route)("users"),
    (0, tsoa_1.Tags)("Users")
], UsersController);
