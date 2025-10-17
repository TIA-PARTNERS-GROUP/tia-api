import { z } from 'zod';

export const createSkillSchema = z.object({
  category: z.string()
    .min(1, 'Category is required')
    .max(100, 'Category must not exceed 100 characters')
    .trim(),
  name: z.string()
    .min(1, 'Skill name is required')
    .max(100, 'Skill name must not exceed 100 characters')
    .trim(),
  description: z.string()
    .max(1000, 'Description must not exceed 1000 characters')
    .optional()
    .default(''),
  active: z.number()
    .int()
    .min(0)
    .max(1)
    .default(1)
});

export const updateSkillSchema = z.object({
  category: z.string()
    .min(1, 'Category is required')
    .max(100, 'Category must not exceed 100 characters')
    .trim()
    .optional(),
  name: z.string()
    .min(1, 'Skill name is required')
    .max(100, 'Skill name must not exceed 100 characters')
    .trim()
    .optional(),
  description: z.string()
    .max(1000, 'Description must not exceed 1000 characters')
    .optional()
    .nullable(),
  active: z.number()
    .int()
    .min(0)
    .max(1)
    .optional()
}).refine(data => Object.keys(data).length > 0, {
  message: 'At least one field must be provided for update'
});

export const skillsFilterSchema = z.object({
  category: z.string().optional(),
  active: z.string().transform(val => val === 'true').optional(),
  search: z.string().optional()
});

export const skillIdSchema = z.object({
  skillId: z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
    message: 'Skill ID must be a positive integer'
  })
});

export const skillNameSchema = z.object({
  name: z.string().min(1, 'Skill name is required').trim()
});

export const skillCategorySchema = z.object({
  category: z.string().min(1, 'Category is required').trim()
});

export const searchSkillsSchema = z.object({
  query: z.string().min(1, 'Search query is required').trim(),
  limit: z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
    message: 'Limit must be a positive integer'
  }).optional()
});

export const popularSkillsSchema = z.object({
  limit: z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
    message: 'Limit must be a positive integer'
  }).optional()
});

export type CreateSkillInput = z.infer<typeof createSkillSchema>;
export type UpdateSkillInput = z.infer<typeof updateSkillSchema>;
export type SkillsFilterInput = z.infer<typeof skillsFilterSchema>;
export type SkillIdParams = z.infer<typeof skillIdSchema>;
export type SkillNameParams = z.infer<typeof skillNameSchema>;
export type SkillCategoryParams = z.infer<typeof skillCategorySchema>;
export type SearchSkillsInput = z.infer<typeof searchSkillsSchema>;
export type PopularSkillsInput = z.infer<typeof popularSkillsSchema>;
