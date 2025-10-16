"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.userUpdateSchema = exports.userCreateSchema = exports.passwordStrengthSchema = exports.loginSchema = void 0;
const zod_1 = require("zod");
exports.loginSchema = zod_1.z.object({
    login_email: zod_1.z.string().email('Please provide a valid email address'),
    password: zod_1.z.string().min(1, 'Password is required')
});
exports.passwordStrengthSchema = zod_1.z.string()
    .min(8, 'Password must be at least 8 characters long')
    .regex(/^(?=.*[a-z])/, 'Password must contain at least one lowercase letter')
    .regex(/^(?=.*[A-Z])/, 'Password must contain at least one uppercase letter')
    .regex(/^(?=.*\d)/, 'Password must contain at least one number')
    .regex(/^(?=.*[@$!%*?&])/, 'Password must contain at least one special character (@$!%*?&)');
exports.userCreateSchema = zod_1.z.object({
    first_name: zod_1.z.string().min(1, 'First name is required'),
    last_name: zod_1.z.string().nullable().optional(),
    login_email: zod_1.z.string().email('Please provide a valid email address'),
    contact_email: zod_1.z.string().email('Please provide a valid contact email').nullable().optional(),
    password: exports.passwordStrengthSchema,
    active: zod_1.z.boolean().default(true),
    email_verified: zod_1.z.boolean().default(false)
});
exports.userUpdateSchema = zod_1.z.object({
    first_name: zod_1.z.string().optional(),
    last_name: zod_1.z.string().nullable().optional(),
    contact_email: zod_1.z.string().email().nullable().optional(),
    password: exports.passwordStrengthSchema.optional(),
    active: zod_1.z.boolean().optional(),
    email_verified: zod_1.z.boolean().optional()
});
