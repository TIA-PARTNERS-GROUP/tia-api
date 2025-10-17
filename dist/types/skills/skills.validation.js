"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.popularSkillsSchema = exports.searchSkillsSchema = exports.skillCategorySchema = exports.skillNameSchema = exports.skillIdSchema = exports.skillsFilterSchema = exports.updateSkillSchema = exports.createSkillSchema = void 0;
const zod_1 = require("zod");
exports.createSkillSchema = zod_1.z.object({
    category: zod_1.z.string()
        .min(1, 'Category is required')
        .max(100, 'Category must not exceed 100 characters')
        .trim(),
    name: zod_1.z.string()
        .min(1, 'Skill name is required')
        .max(100, 'Skill name must not exceed 100 characters')
        .trim(),
    description: zod_1.z.string()
        .max(1000, 'Description must not exceed 1000 characters')
        .optional()
        .default(''),
    active: zod_1.z.number()
        .int()
        .min(0)
        .max(1)
        .default(1)
});
exports.updateSkillSchema = zod_1.z.object({
    category: zod_1.z.string()
        .min(1, 'Category is required')
        .max(100, 'Category must not exceed 100 characters')
        .trim()
        .optional(),
    name: zod_1.z.string()
        .min(1, 'Skill name is required')
        .max(100, 'Skill name must not exceed 100 characters')
        .trim()
        .optional(),
    description: zod_1.z.string()
        .max(1000, 'Description must not exceed 1000 characters')
        .optional()
        .nullable(),
    active: zod_1.z.number()
        .int()
        .min(0)
        .max(1)
        .optional()
}).refine(data => Object.keys(data).length > 0, {
    message: 'At least one field must be provided for update'
});
exports.skillsFilterSchema = zod_1.z.object({
    category: zod_1.z.string().optional(),
    active: zod_1.z.string().transform(val => val === 'true').optional(),
    search: zod_1.z.string().optional()
});
exports.skillIdSchema = zod_1.z.object({
    skillId: zod_1.z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
        message: 'Skill ID must be a positive integer'
    })
});
exports.skillNameSchema = zod_1.z.object({
    name: zod_1.z.string().min(1, 'Skill name is required').trim()
});
exports.skillCategorySchema = zod_1.z.object({
    category: zod_1.z.string().min(1, 'Category is required').trim()
});
exports.searchSkillsSchema = zod_1.z.object({
    query: zod_1.z.string().min(1, 'Search query is required').trim(),
    limit: zod_1.z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
        message: 'Limit must be a positive integer'
    }).optional()
});
exports.popularSkillsSchema = zod_1.z.object({
    limit: zod_1.z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
        message: 'Limit must be a positive integer'
    }).optional()
});
