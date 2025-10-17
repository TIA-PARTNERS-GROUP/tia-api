"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.passwordStrengthSchema = exports.loginSchema = void 0;
const zod_1 = require("zod");
exports.loginSchema = zod_1.z.object({
    login_email: zod_1.z.string().email("Please provide a valid email address"),
    password: zod_1.z.string().min(1, "Password is required"),
});
exports.passwordStrengthSchema = zod_1.z
    .string()
    .min(8, "Password must be at least 8 characters long")
    .regex(/^(?=.*[a-z])/, "Password must contain at least one lowercase letter")
    .regex(/^(?=.*[A-Z])/, "Password must contain at least one uppercase letter")
    .regex(/^(?=.*\d)/, "Password must contain at least one number")
    .regex(/^(?=.*[@$!%*?&])/, "Password must contain at least one special character (@$!%*?&)");
