"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.proficiencyLevelSchema = exports.skillIdParamsSchema = exports.userIdSchema = exports.userSkillIdSchema = exports.updateUserSkillSchema = exports.createUserSkillSchema = void 0;
const zod_1 = require("zod");
const proficiencyLevelEnum = zod_1.z.enum(['beginner', 'intermediate', 'advanced', 'expert']);
exports.createUserSkillSchema = zod_1.z.object({
    skill_id: zod_1.z.number()
        .int()
        .positive('Skill ID must be a positive integer'),
    user_id: zod_1.z.number()
        .int()
        .positive('User ID must be a positive integer'),
    proficiency_level: proficiencyLevelEnum
        .default('intermediate')
});
exports.updateUserSkillSchema = zod_1.z.object({
    proficiency_level: proficiencyLevelEnum
});
exports.userSkillIdSchema = zod_1.z.object({
    skillId: zod_1.z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
        message: 'Skill ID must be a positive integer'
    }),
    userId: zod_1.z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
        message: 'User ID must be a positive integer'
    })
});
exports.userIdSchema = zod_1.z.object({
    userId: zod_1.z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
        message: 'User ID must be a positive integer'
    })
});
exports.skillIdParamsSchema = zod_1.z.object({
    skillId: zod_1.z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
        message: 'Skill ID must be a positive integer'
    })
});
exports.proficiencyLevelSchema = zod_1.z.object({
    userId: zod_1.z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
        message: 'User ID must be a positive integer'
    }),
    proficiencyLevel: proficiencyLevelEnum
});
