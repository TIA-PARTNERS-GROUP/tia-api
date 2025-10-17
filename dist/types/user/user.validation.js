"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.userIdParamsSchema = exports.userUpdateSchema = exports.userCreateSchema = void 0;
const zod_1 = require("zod");
const password_utils_1 = require("../../utils/password.utils");
const passwordSchema = zod_1.z
    .string()
    .min(8, "Password must be at least 8 characters long")
    .max(100, "Password too long")
    .refine((password) => {
    const validation = password_utils_1.PasswordUtils.validatePasswordComplexity(password);
    return validation.isValid;
}, {
    message: "Password must contain at least one lowercase letter, one uppercase letter, one number, and one special character",
});
exports.userCreateSchema = zod_1.z
    .object({
    first_name: zod_1.z
        .string()
        .trim()
        .min(1, "First name is required")
        .max(60, "First name must be less than 60 characters"),
    last_name: zod_1.z
        .string()
        .trim()
        .min(1, "Last name is required")
        .max(60, "Last name must be less than 60 characters"),
    login_email: zod_1.z
        .string()
        .trim()
        .min(1, "Login email is required")
        .email("Invalid email address")
        .max(254, "Email must be less than 254 characters"),
    password: passwordSchema,
    contact_email: zod_1.z
        .string()
        .email("Invalid contact email address")
        .max(254, "Contact email must be less than 254 characters")
        .optional()
        .transform((val) => val ?? null),
    contact_phone_no: zod_1.z
        .string()
        .max(20, "Phone number must be less than 20 characters")
        .optional()
        .transform((val) => val ?? null),
    adk_session_id: zod_1.z
        .string()
        .uuid("Must be a valid UUID")
        .max(128, "Session ID too long")
        .optional()
        .transform((val) => val ?? null),
})
    .strict();
exports.userUpdateSchema = exports.userCreateSchema
    .partial()
    .extend({
    email_verified: zod_1.z.boolean().optional(),
    active: zod_1.z.boolean().optional(),
})
    .refine((data) => {
    if (data.password) {
        const validation = password_utils_1.PasswordUtils.validatePasswordComplexity(data.password);
        return validation.isValid;
    }
    return true;
}, {
    message: "Password must contain at least one lowercase letter, one uppercase letter, one number, and one special character",
    path: ["password"],
});
exports.userIdParamsSchema = zod_1.z.object({
    id: zod_1.z
        .string()
        .transform((val) => parseInt(val, 10))
        .refine((val) => !isNaN(val) && val > 0, {
        message: "User ID must be a positive integer",
    }),
});
