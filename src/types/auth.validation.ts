import { z } from 'zod';

export const loginSchema = z.object({
  login_email: z.string().email('Please provide a valid email address'),
  password: z.string().min(1, 'Password is required')
});

export const passwordStrengthSchema = z.string()
  .min(8, 'Password must be at least 8 characters long')
  .regex(/^(?=.*[a-z])/, 'Password must contain at least one lowercase letter')
  .regex(/^(?=.*[A-Z])/, 'Password must contain at least one uppercase letter')
  .regex(/^(?=.*\d)/, 'Password must contain at least one number')
  .regex(/^(?=.*[@$!%*?&])/, 'Password must contain at least one special character (@$!%*?&)');

export const userCreateSchema = z.object({
  first_name: z.string().min(1, 'First name is required'),
  last_name: z.string().nullable().optional(),
  login_email: z.string().email('Please provide a valid email address'),
  contact_email: z.string().email('Please provide a valid contact email').nullable().optional(),
  password: passwordStrengthSchema,
  active: z.boolean().default(true),
  email_verified: z.boolean().default(false)
});

export const userUpdateSchema = z.object({
  first_name: z.string().optional(),
  last_name: z.string().nullable().optional(),
  contact_email: z.string().email().nullable().optional(),
  password: passwordStrengthSchema.optional(),
  active: z.boolean().optional(),
  email_verified: z.boolean().optional()
});

export type LoginInput = z.infer<typeof loginSchema>;
export type UserCreateInput = z.infer<typeof userCreateSchema>;
export type UserUpdateInput = z.infer<typeof userUpdateSchema>;
