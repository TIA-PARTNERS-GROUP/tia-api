import { prisma } from '@lib/prisma';
import { ZodError } from 'zod';
import { StatusCodes } from 'http-status-codes';
import { ApiError } from '../errors/ApiError';
import { createUserSchema, updateUserSchema } from '../types/user.validation';
import type { UserCreationParams, UserUpdateParams } from '../types/user.validation';
import type { UserResponse, UserAuthResponse, LoginRequest, LoginResponse } from '../types/user.dto';
import { PasswordUtils } from '../utils/password.utils';

function isPrismaError(error: unknown): error is { code: string } {
  return typeof error === 'object' && error !== null && 'code' in error;
}

export const createUser = async (data: UserCreationParams): Promise<UserResponse> => {
  try {
    const validatedData = createUserSchema.parse(data);

   
    const hashedPassword = await PasswordUtils.hashPassword(validatedData.password);

    const user = await prisma.users.create({
      data: {
        first_name: validatedData.first_name,
        last_name: validatedData.last_name,
        login_email: validatedData.login_email,
        password_hash: hashedPassword,
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
  } catch (error) {
    if (error instanceof ZodError) {
      const errorDetails = error.issues.map((issue) => ({
        field: issue.path.join('.'),
        message: issue.message,
      }));
      throw new ApiError(StatusCodes.UNPROCESSABLE_ENTITY, 'Validation failed', errorDetails);
    }

   
    if (isPrismaError(error)) {
      if (error.code === 'P2002') {
        throw new ApiError(StatusCodes.CONFLICT, 'User with this email already exists');
      }
    }

    throw new ApiError(StatusCodes.INTERNAL_SERVER_ERROR, 'An unexpected error occurred');
  }
};

export const updateUser = async (id: number, data: UserUpdateParams): Promise<UserResponse | null> => {
  try {
    const validatedData = updateUserSchema.parse(data);

   
    const existingUser = await prisma.users.findUnique({ where: { id } });
    if (!existingUser) {
      return null;
    }

   
    const updateData: any = {
      first_name: validatedData.first_name,
      last_name: validatedData.last_name,
      login_email: validatedData.login_email,
      contact_email: validatedData.contact_email,
      contact_phone_no: validatedData.contact_phone_no,
      adk_session_id: validatedData.adk_session_id,
      email_verified: validatedData.email_verified,
      active: validatedData.active,
    };

   
    if (validatedData.password) {
      updateData.password_hash = await PasswordUtils.hashPassword(validatedData.password);
    }

   
    Object.keys(updateData).forEach(key => {
      if (updateData[key] === undefined) {
        delete updateData[key];
      }
    });

    if (Object.keys(updateData).length === 0) {
      throw new ApiError(StatusCodes.BAD_REQUEST, 'No valid fields provided for update.');
    }

    const user = await prisma.users.update({
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
  } catch (error) {
    if (error instanceof ZodError) {
      const errorDetails = error.issues.map((issue) => ({
        field: issue.path.join('.'),
        message: issue.message,
      }));
      throw new ApiError(StatusCodes.UNPROCESSABLE_ENTITY, 'Validation failed', errorDetails);
    }

   
    if (isPrismaError(error)) {
      if (error.code === 'P2002') {
        throw new ApiError(StatusCodes.CONFLICT, 'User with this email already exists');
      }
    }

    throw error;
  }
};

export const authenticateUser = async (loginData: LoginRequest): Promise<LoginResponse> => {
  try {
    const { login_email, password } = loginData;

   
    const user = await prisma.users.findUnique({
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
      throw new ApiError(StatusCodes.UNAUTHORIZED, 'Invalid email or password');
    }

    if (!user.active) {
      throw new ApiError(StatusCodes.UNAUTHORIZED, 'Account is deactivated');
    }

    if (!user.password_hash) {
      throw new ApiError(StatusCodes.UNAUTHORIZED, 'No password set for this account');
    }

   
    const isPasswordValid = await PasswordUtils.verifyPassword(password, user.password_hash);
    if (!isPasswordValid) {
      throw new ApiError(StatusCodes.UNAUTHORIZED, 'Invalid email or password');
    }

   
    const { password_hash, ...userWithoutPassword } = user;

    return {
      user: userWithoutPassword,
     
    };
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    throw new ApiError(StatusCodes.INTERNAL_SERVER_ERROR, 'Authentication failed');
  }
};

export const changePassword = async (userId: number, currentPassword: string, newPassword: string): Promise<boolean> => {
  try {
    const user = await prisma.users.findUnique({
      where: { id: userId },
      select: { password_hash: true }
    });

    if (!user || !user.password_hash) {
      throw new ApiError(StatusCodes.NOT_FOUND, 'User not found');
    }

   
    const isCurrentPasswordValid = await PasswordUtils.verifyPassword(currentPassword, user.password_hash);
    if (!isCurrentPasswordValid) {
      throw new ApiError(StatusCodes.UNAUTHORIZED, 'Current password is incorrect');
    }

   
    const complexityCheck = PasswordUtils.validatePasswordComplexity(newPassword);
    if (!complexityCheck.isValid) {
      throw new ApiError(StatusCodes.UNPROCESSABLE_ENTITY, complexityCheck.message || 'Password does not meet requirements');
    }

   
    const newHashedPassword = await PasswordUtils.hashPassword(newPassword);
    await prisma.users.update({
      where: { id: userId },
      data: { password_hash: newHashedPassword }
    });

    return true;
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    throw new ApiError(StatusCodes.INTERNAL_SERVER_ERROR, 'Failed to change password');
  }
};

export const findUserById = async (id: number): Promise<UserResponse | null> => {
  const user = await prisma.users.findUnique({
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

export const findAllUsers = async (): Promise<UserResponse[]> => {
  const users = await prisma.users.findMany({
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

export const deleteUser = async (id: number): Promise<UserResponse | null> => {
  try {
    const user = await prisma.users.delete({
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
  } catch (error) {
   
    if (isPrismaError(error)) {
      if (error.code === 'P2025') {
        return null;
      }
    }
    throw error;
  }
};

export const findUserByEmail = async (email: string): Promise<UserAuthResponse | null> => {
  const user = await prisma.users.findUnique({
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

export const deactivateUser = async (id: number): Promise<UserResponse | null> => {
  const user = await prisma.users.update({
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
