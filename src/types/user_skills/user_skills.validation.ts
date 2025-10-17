import { z } from 'zod';

const proficiencyLevelEnum = z.enum(['beginner', 'intermediate', 'advanced', 'expert']);

export const createUserSkillSchema = z.object({
  skill_id: z.number()
    .int()
    .positive('Skill ID must be a positive integer'),
  user_id: z.number()
    .int()
    .positive('User ID must be a positive integer'),
  proficiency_level: proficiencyLevelEnum
    .default('intermediate')
});

export const updateUserSkillSchema = z.object({
  proficiency_level: proficiencyLevelEnum
});

export const userSkillIdSchema = z.object({
  skillId: z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
    message: 'Skill ID must be a positive integer'
  }),
  userId: z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
    message: 'User ID must be a positive integer'
  })
});

export const userIdSchema = z.object({
  userId: z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
    message: 'User ID must be a positive integer'
  })
});

export const skillIdParamsSchema = z.object({
  skillId: z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
    message: 'Skill ID must be a positive integer'
  })
});

export const proficiencyLevelSchema = z.object({
  userId: z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
    message: 'User ID must be a positive integer'
  }),
  proficiencyLevel: proficiencyLevelEnum
});

export type CreateUserSkillInput = z.infer<typeof createUserSkillSchema>;
export type UpdateUserSkillInput = z.infer<typeof updateUserSkillSchema>;
export type UserSkillIdParams = z.infer<typeof userSkillIdSchema>;
export type UserIdParams = z.infer<typeof userIdSchema>;
export type SkillIdParams = z.infer<typeof skillIdParamsSchema>;
export type ProficiencyLevelParams = z.infer<typeof proficiencyLevelSchema>;
