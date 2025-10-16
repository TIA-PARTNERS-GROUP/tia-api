"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivateUser = exports.findUserByEmail = exports.deleteUser = exports.findAllUsers = exports.findUserById = exports.changePassword = exports.authenticateUser = exports.updateUser = exports.createUser = void 0;
const prisma_1 = require("../lib/prisma");
const zod_1 = require("zod");
const http_status_codes_1 = require("http-status-codes");
const ApiError_1 = require("../errors/ApiError");
const user_validation_1 = require("../types/user.validation");
const password_utils_1 = require("../utils/password.utils");
function isPrismaError(error) {
    return typeof error === 'object' && error !== null && 'code' in error;
}
const createUser = async (data) => {
    try {
        const validatedData = user_validation_1.createUserSchema.parse(data);
        // Hash the password before storing
        const hashedPassword = await password_utils_1.PasswordUtils.hashPassword(validatedData.password);
        const user = await prisma_1.prisma.users.create({
            data: {
                first_name: validatedData.first_name,
                last_name: validatedData.last_name,
                login_email: validatedData.login_email,
                password_hash: hashedPassword, // Store the hashed password
                contact_email: validatedData.contact_email,
                contact_phone_no: validatedData.contact_phone_no,
                adk_session_id: validatedData.adk_session_id,
                email_verified: false,
                active: true,
            },
            select: {
                id: true,
                first_name: true,
                last_name: true,
                login_email: true,
                contact_email: true,
                contact_phone_no: true,
                adk_session_id: true,
                email_verified: true,
                active: true,
                created_at: true,
                updated_at: true,
            }
        });
        return user;
    }
    catch (error) {
        if (error instanceof zod_1.ZodError) {
            const errorDetails = error.issues.map((issue) => ({
                field: issue.path.join('.'),
                message: issue.message,
            }));
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.UNPROCESSABLE_ENTITY, 'Validation failed', errorDetails);
        }
        // Safe Prisma error handling
        if (isPrismaError(error)) {
            if (error.code === 'P2002') {
                throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.CONFLICT, 'User with this email already exists');
            }
        }
        throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR, 'An unexpected error occurred');
    }
};
exports.createUser = createUser;
const updateUser = async (id, data) => {
    try {
        const validatedData = user_validation_1.updateUserSchema.parse(data);
        // Check if user exists
        const existingUser = await prisma_1.prisma.users.findUnique({ where: { id } });
        if (!existingUser) {
            return null;
        }
        // Prepare update data
        const updateData = {
            first_name: validatedData.first_name,
            last_name: validatedData.last_name,
            login_email: validatedData.login_email,
            contact_email: validatedData.contact_email,
            contact_phone_no: validatedData.contact_phone_no,
            adk_session_id: validatedData.adk_session_id,
            email_verified: validatedData.email_verified,
            active: validatedData.active,
        };
        // Hash password if it's being updated
        if (validatedData.password) {
            updateData.password_hash = await password_utils_1.PasswordUtils.hashPassword(validatedData.password);
        }
        // Filter out undefined values
        Object.keys(updateData).forEach(key => {
            if (updateData[key] === undefined) {
                delete updateData[key];
            }
        });
        if (Object.keys(updateData).length === 0) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.BAD_REQUEST, 'No valid fields provided for update.');
        }
        const user = await prisma_1.prisma.users.update({
            where: { id },
            data: updateData,
            select: {
                id: true,
                first_name: true,
                last_name: true,
                login_email: true,
                contact_email: true,
                contact_phone_no: true,
                adk_session_id: true,
                email_verified: true,
                active: true,
                created_at: true,
                updated_at: true,
            }
        });
        return user;
    }
    catch (error) {
        if (error instanceof zod_1.ZodError) {
            const errorDetails = error.issues.map((issue) => ({
                field: issue.path.join('.'),
                message: issue.message,
            }));
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.UNPROCESSABLE_ENTITY, 'Validation failed', errorDetails);
        }
        // Safe Prisma error handling
        if (isPrismaError(error)) {
            if (error.code === 'P2002') {
                throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.CONFLICT, 'User with this email already exists');
            }
        }
        throw error;
    }
};
exports.updateUser = updateUser;
// Authentication service functions
const authenticateUser = async (loginData) => {
    try {
        const { login_email, password } = loginData;
        // Find user by email
        const user = await prisma_1.prisma.users.findUnique({
            where: { login_email },
            select: {
                id: true,
                first_name: true,
                last_name: true,
                login_email: true,
                password_hash: true,
                contact_email: true,
                email_verified: true,
                active: true,
                created_at: true,
            }
        });
        if (!user) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Invalid email or password');
        }
        if (!user.active) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Account is deactivated');
        }
        if (!user.password_hash) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'No password set for this account');
        }
        // Verify password
        const isPasswordValid = await password_utils_1.PasswordUtils.verifyPassword(password, user.password_hash);
        if (!isPasswordValid) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Invalid email or password');
        }
        // Remove password hash from response
        const { password_hash, ...userWithoutPassword } = user;
        return {
            user: userWithoutPassword,
            // You can add JWT token or session ID here later
        };
    }
    catch (error) {
        if (error instanceof ApiError_1.ApiError) {
            throw error;
        }
        throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR, 'Authentication failed');
    }
};
exports.authenticateUser = authenticateUser;
const changePassword = async (userId, currentPassword, newPassword) => {
    try {
        const user = await prisma_1.prisma.users.findUnique({
            where: { id: userId },
            select: { password_hash: true }
        });
        if (!user || !user.password_hash) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'User not found');
        }
        // Verify current password
        const isCurrentPasswordValid = await password_utils_1.PasswordUtils.verifyPassword(currentPassword, user.password_hash);
        if (!isCurrentPasswordValid) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.UNAUTHORIZED, 'Current password is incorrect');
        }
        // Validate new password complexity
        const complexityCheck = password_utils_1.PasswordUtils.validatePasswordComplexity(newPassword);
        if (!complexityCheck.isValid) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.UNPROCESSABLE_ENTITY, complexityCheck.message || 'Password does not meet requirements');
        }
        // Hash and update new password
        const newHashedPassword = await password_utils_1.PasswordUtils.hashPassword(newPassword);
        await prisma_1.prisma.users.update({
            where: { id: userId },
            data: { password_hash: newHashedPassword }
        });
        return true;
    }
    catch (error) {
        if (error instanceof ApiError_1.ApiError) {
            throw error;
        }
        throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR, 'Failed to change password');
    }
};
exports.changePassword = changePassword;
const findUserById = async (id) => {
    const user = await prisma_1.prisma.users.findUnique({
        where: { id },
        select: {
            id: true,
            first_name: true,
            last_name: true,
            login_email: true,
            contact_email: true,
            contact_phone_no: true,
            adk_session_id: true,
            email_verified: true,
            active: true,
            created_at: true,
            updated_at: true,
        }
    });
    return user;
};
exports.findUserById = findUserById;
const findAllUsers = async () => {
    const users = await prisma_1.prisma.users.findMany({
        select: {
            id: true,
            first_name: true,
            last_name: true,
            login_email: true,
            contact_email: true,
            contact_phone_no: true,
            adk_session_id: true,
            email_verified: true,
            active: true,
            created_at: true,
            updated_at: true,
        },
        orderBy: {
            created_at: 'desc'
        }
    });
    return users;
};
exports.findAllUsers = findAllUsers;
const deleteUser = async (id) => {
    try {
        const user = await prisma_1.prisma.users.delete({
            where: { id },
            select: {
                id: true,
                first_name: true,
                last_name: true,
                login_email: true,
                contact_email: true,
                contact_phone_no: true,
                adk_session_id: true,
                email_verified: true,
                active: true,
                created_at: true,
                updated_at: true,
            }
        });
        return user;
    }
    catch (error) {
        // Safe Prisma error handling
        if (isPrismaError(error)) {
            if (error.code === 'P2025') {
                return null; // User not found
            }
        }
        throw error;
    }
};
exports.deleteUser = deleteUser;
const findUserByEmail = async (email) => {
    const user = await prisma_1.prisma.users.findUnique({
        where: { login_email: email },
        select: {
            id: true,
            first_name: true,
            last_name: true,
            login_email: true,
            password_hash: true,
            contact_email: true,
            email_verified: true,
            active: true,
            created_at: true,
        }
    });
    return user;
};
exports.findUserByEmail = findUserByEmail;
const deactivateUser = async (id) => {
    const user = await prisma_1.prisma.users.update({
        where: { id },
        data: { active: false },
        select: {
            id: true,
            first_name: true,
            last_name: true,
            login_email: true,
            contact_email: true,
            contact_phone_no: true,
            adk_session_id: true,
            email_verified: true,
            active: true,
            created_at: true,
            updated_at: true,
        }
    });
    return user;
};
exports.deactivateUser = deactivateUser;
