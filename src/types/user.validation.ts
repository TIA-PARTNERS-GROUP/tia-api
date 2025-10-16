import { z } from 'zod';
import { PasswordUtils } from '../utils/password.utils';

export interface UserCreationParams {
  first_name: string;
  last_name: string;
  login_email: string;
  password: string;
  contact_email?: string;
  contact_phone_no?: string;
  adk_session_id?: string;
}

export interface UserUpdateParams {
  first_name?: string | undefined;
  last_name?: string | undefined;
  login_email?: string | undefined;
  password?: string | undefined;
  contact_email?: string | undefined;
  contact_phone_no?: string | undefined;
  adk_session_id?: string | undefined;
  email_verified?: boolean | undefined;
  active?: boolean | undefined;
}

const passwordSchema = z.string()
  .min(8, 'Password must be at least 8 characters long')
  .max(100, 'Password too long')
  .refine(
    (password) => {
      const validation = PasswordUtils.validatePasswordComplexity(password);
      return validation.isValid;
    },
    {
      message: 'Password must contain at least one lowercase letter, one uppercase letter, one number, and one special character'
    }
  );

export const createUserSchema = z.object({
  first_name: z
    .string()
    .trim()
    .min(1, 'First name is required')
    .max(60, 'First name must be less than 60 characters'),
  last_name: z
    .string()
    .trim()
    .min(1, 'Last name is required')
    .max(60, 'Last name must be less than 60 characters'),
  login_email: z
    .string()
    .trim()
    .min(1, 'Login email is required')
    .email('Invalid email address')
    .max(254, 'Email must be less than 254 characters'),
  password: passwordSchema,
  contact_email: z
    .string()
    .email('Invalid contact email address')
    .max(254, 'Contact email must be less than 254 characters')
    .optional()
    .transform(val => val ?? null),
  contact_phone_no: z
    .string()
    .max(20, 'Phone number must be less than 20 characters')
    .optional()
    .transform(val => val ?? null),
  adk_session_id: z
    .string()
    .uuid('Must be a valid UUID')
    .max(128, 'Session ID too long')
    .optional()
    .transform(val => val ?? null),
}).strict();

export const updateUserSchema = createUserSchema.partial().extend({
  email_verified: z.boolean().optional(),
  active: z.boolean().optional(),
}).refine(
  (data) => {
    if (data.password) {
      const validation = PasswordUtils.validatePasswordComplexity(data.password);
      return validation.isValid;
    }
    return true;
  },
  {
    message: 'Password must contain at least one lowercase letter, one uppercase letter, one number, and one special character',
    path: ['password']
  }
);
