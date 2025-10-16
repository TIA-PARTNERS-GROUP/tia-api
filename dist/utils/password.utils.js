"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.PasswordUtils = void 0;
const bcryptjs_1 = __importDefault(require("bcryptjs"));
class PasswordUtils {
    static SALT_ROUNDS = 12;
    static async hashPassword(plainPassword) {
        return await bcryptjs_1.default.hash(plainPassword, this.SALT_ROUNDS);
    }
    static async verifyPassword(plainPassword, hashedPassword) {
        return await bcryptjs_1.default.compare(plainPassword, hashedPassword);
    }
    static validatePasswordComplexity(password) {
        if (password.length < 8) {
            return { isValid: false, message: 'Password must be at least 8 characters long' };
        }
        if (!/(?=.*[a-z])/.test(password)) {
            return { isValid: false, message: 'Password must contain at least one lowercase letter' };
        }
        if (!/(?=.*[A-Z])/.test(password)) {
            return { isValid: false, message: 'Password must contain at least one uppercase letter' };
        }
        if (!/(?=.*\d)/.test(password)) {
            return { isValid: false, message: 'Password must contain at least one number' };
        }
        if (!/(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?])/.test(password)) {
            return { isValid: false, message: 'Password must contain at least one special character' };
        }
        return { isValid: true };
    }
}
exports.PasswordUtils = PasswordUtils;
