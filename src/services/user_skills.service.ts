import { prisma } from '@lib/prisma';
import { createUserSkillSchema, updateUserSkillSchema, CreateUserSkillInput, UpdateUserSkillInput } from '../types/user_skills.validation';
import { ApiError } from '../errors/ApiError';
import { StatusCodes } from 'http-status-codes';

export class UserSkillsService {
  /**
   * Get all skills for a user
   */
  async getUserSkills(userId: number) {
    return await prisma.user_skills.findMany({
      where: { user_id: userId },
      include: {
        skills: {
          select: {
            id: true,
            category: true,
            name: true,
            description: true
          }
        },
        users: {
          select: {
            id: true,
            first_name: true,
            last_name: true,
            login_email: true
          }
        }
      }
    });
  }

  /**
   * Get a specific user skill by skill_id and user_id
   */
  async getUserSkillById(skillId: number, userId: number) {
    const userSkill = await prisma.user_skills.findUnique({
      where: {
        skill_id_user_id: {
          skill_id: skillId,
          user_id: userId
        }
      },
      include: {
        skills: {
          select: {
            id: true,
            category: true,
            name: true,
            description: true
          }
        },
        users: {
          select: {
            id: true,
            first_name: true,
            last_name: true,
            login_email: true
          }
        }
      }
    });

    if (!userSkill) {
      throw new ApiError(StatusCodes.NOT_FOUND, 'User skill not found');
    }

    return userSkill;
  }

  /**
   * Create a new user skill
   */
  async createUserSkill(data: any) {
   
    const validatedData = createUserSkillSchema.parse(data);

   
    const user = await prisma.users.findUnique({
      where: { id: validatedData.user_id }
    });

    if (!user) {
      throw new ApiError(StatusCodes.NOT_FOUND, 'User not found');
    }

   
    const skill = await prisma.skills.findUnique({
      where: { id: validatedData.skill_id }
    });

    if (!skill) {
      throw new ApiError(StatusCodes.NOT_FOUND, 'Skill not found');
    }

   
    const existingUserSkill = await prisma.user_skills.findUnique({
      where: {
        skill_id_user_id: {
          skill_id: validatedData.skill_id,
          user_id: validatedData.user_id
        }
      }
    });

    if (existingUserSkill) {
      throw new ApiError(StatusCodes.CONFLICT, 'User already has this skill');
    }

   
    return await prisma.user_skills.create({
      data: {
        skill_id: validatedData.skill_id,
        user_id: validatedData.user_id,
        proficiency_level: validatedData.proficiency_level
      },
      include: {
        skills: {
          select: {
            id: true,
            category: true,
            name: true,
            description: true
          }
        },
        users: {
          select: {
            id: true,
            first_name: true,
            last_name: true,
            login_email: true
          }
        }
      }
    });
  }
  /**
   * Update a user skill
   */

  async updateUserSkill(skillId: number, userId: number, data: any) {
   
    const validatedData = updateUserSkillSchema.parse(data);

   
    const existingUserSkill = await prisma.user_skills.findUnique({
      where: {
        skill_id_user_id: {
          skill_id: skillId,
          user_id: userId
        }
      }
    });

    if (!existingUserSkill) {
      throw new ApiError(StatusCodes.NOT_FOUND, 'User skill not found');
    }

   
    return await prisma.user_skills.update({
      where: {
        skill_id_user_id: {
          skill_id: skillId,
          user_id: userId
        }
      },
      data: {
        proficiency_level: validatedData.proficiency_level
      },
      include: {
        skills: {
          select: {
            id: true,
            category: true,
            name: true,
            description: true
          }
        },
        users: {
          select: {
            id: true,
            first_name: true,
            last_name: true,
            login_email: true
          }
        }
      }
    });
  }

  /**
   * Delete a user skill
   */
  async deleteUserSkill(skillId: number, userId: number) {
    const existingUserSkill = await prisma.user_skills.findUnique({
      where: {
        skill_id_user_id: {
          skill_id: skillId,
          user_id: userId
        }
      }
    });

    if (!existingUserSkill) {
      throw new ApiError(StatusCodes.NOT_FOUND, 'User skill not found');
    }

    await prisma.user_skills.delete({
      where: {
        skill_id_user_id: {
          skill_id: skillId,
          user_id: userId
        }
      }
    });

    return { message: 'User skill deleted successfully' };
  }

  /**
   * Get users by skill
   */
  async getUsersBySkill(skillId: number) {
    const skill = await prisma.skills.findUnique({
      where: { id: skillId }
    });

    if (!skill) {
      throw new ApiError(StatusCodes.NOT_FOUND, 'Skill not found');
    }

    return await prisma.user_skills.findMany({
      where: { skill_id: skillId },
      include: {
        users: {
          select: {
            id: true,
            first_name: true,
            last_name: true,
            login_email: true,
            contact_email: true
          }
        },
        skills: {
          select: {
            id: true,
            category: true,
            name: true,
            description: true
          }
        }
      }
    });
  }

  /**
   * Get skills by proficiency level for a user
   */
  async getUserSkillsByProficiency(userId: number, proficiencyLevel: string) {
    return await prisma.user_skills.findMany({
      where: {
        user_id: userId,
        proficiency_level: proficiencyLevel as any
      },
      include: {
        skills: {
          select: {
            id: true,
            category: true,
            name: true,
            description: true
          }
        }
      }
    });
  }
}

export const userSkillsService = new UserSkillsService();
