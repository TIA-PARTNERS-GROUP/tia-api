import { prisma } from '@lib/prisma';
import { SkillsFilter } from '../types/skills.types';
import { createSkillSchema, updateSkillSchema } from '../types/skills.validation';
import { ApiError } from '../errors/ApiError';
import { StatusCodes } from 'http-status-codes';

export class SkillsService {
  /**
   * Get all skills with optional filtering
   */
  async getSkills(filters?: SkillsFilter) {
    const where: any = {};

    if (filters?.category) {
      where.category = { contains: filters.category };
    }

    if (filters?.active !== undefined) {
      where.active = filters.active ? 1 : 0;
    }

    if (filters?.search) {
      where.OR = [
        { name: { contains: filters.search } },
        { category: { contains: filters.search } },
        { description: { contains: filters.search } }
      ];
    }

    return await prisma.skills.findMany({
      where,
      include: {
        _count: {
          select: {
            user_skills: true,
            project_skills: true
          }
        }
      },
      orderBy: {
        name: 'asc'
      }
    });
  }

  /**
   * Get a specific skill by ID
   */
  async getSkillById(skillId: number) {
    const skill = await prisma.skills.findUnique({
      where: { id: skillId },
      include: {
        _count: {
          select: {
            user_skills: true,
            project_skills: true
          }
        },
        user_skills: {
          include: {
            users: {
              select: {
                id: true,
                first_name: true,
                last_name: true,
                login_email: true
              }
            }
          },
          take: 10
        },
        project_skills: {
          include: {
            projects: {
              select: {
                id: true,
                name: true,
                project_status: true
              }
            }
          },
          take: 10
        }
      }
    });

    if (!skill) {
      throw new ApiError(StatusCodes.NOT_FOUND, 'Skill not found');
    }

    return skill;
  }

  /**
   * Get a skill by name
   */
  async getSkillByName(name: string) {
    const skill = await prisma.skills.findUnique({
      where: { name },
      include: {
        _count: {
          select: {
            user_skills: true,
            project_skills: true
          }
        }
      }
    });

    if (!skill) {
      throw new ApiError(StatusCodes.NOT_FOUND, 'Skill not found');
    }

    return skill;
  }

  /**
   * Create a new skill
   */
  async createSkill(data: any) {

    const validatedData = createSkillSchema.parse(data);


    const existingSkill = await prisma.skills.findUnique({
      where: { name: validatedData.name }
    });

    if (existingSkill) {
      throw new ApiError(StatusCodes.CONFLICT, 'Skill with this name already exists');
    }


    return await prisma.skills.create({
      data: {
        category: validatedData.category,
        name: validatedData.name,
        description: validatedData.description,
        active: validatedData.active
      },
      include: {
        _count: {
          select: {
            user_skills: true,
            project_skills: true
          }
        }
      }
    });
  }  /**
   * Update a skill
   */
  async updateSkill(skillId: number, data: any) {

    const validatedData = updateSkillSchema.parse(data);


    const existingSkill = await prisma.skills.findUnique({
      where: { id: skillId }
    });

    if (!existingSkill) {
      throw new ApiError(StatusCodes.NOT_FOUND, 'Skill not found');
    }


    if (validatedData.name && validatedData.name !== existingSkill.name) {
      const nameConflict = await prisma.skills.findUnique({
        where: { name: validatedData.name }
      });

      if (nameConflict) {
        throw new ApiError(StatusCodes.CONFLICT, 'Skill with this name already exists');
      }
    }


    const updateData: any = {};
    if (validatedData.category !== undefined) updateData.category = validatedData.category;
    if (validatedData.name !== undefined) updateData.name = validatedData.name;
    if (validatedData.description !== undefined) updateData.description = validatedData.description;
    if (validatedData.active !== undefined) updateData.active = validatedData.active;

    return await prisma.skills.update({
      where: { id: skillId },
      data: updateData,
      include: {
        _count: {
          select: {
            user_skills: true,
            project_skills: true
          }
        }
      }
    });
  }  /**
   * Delete a skill
   */
  async deleteSkill(skillId: number) {

    const existingSkill = await prisma.skills.findUnique({
      where: { id: skillId }
    });

    if (!existingSkill) {
      throw new ApiError(StatusCodes.NOT_FOUND, 'Skill not found');
    }


    const userSkillsCount = await prisma.user_skills.count({
      where: { skill_id: skillId }
    });

    const projectSkillsCount = await prisma.project_skills.count({
      where: { skill_id: skillId }
    });

    if (userSkillsCount > 0 || projectSkillsCount > 0) {
      throw new ApiError(
        StatusCodes.CONFLICT,
        `Cannot delete skill. It is currently used by ${userSkillsCount} users and ${projectSkillsCount} projects.`
      );
    }

    await prisma.skills.delete({
      where: { id: skillId }
    });

    return { message: 'Skill deleted successfully' };
  }

  /**
   * Get skills by category
   */
  async getSkillsByCategory(category: string) {
    return await prisma.skills.findMany({
      where: {
        category: { contains: category },
        active: 1
      },
      include: {
        _count: {
          select: {
            user_skills: true,
            project_skills: true
          }
        }
      },
      orderBy: {
        name: 'asc'
      }
    });
  }

  /**
   * Get all unique categories
   */
  async getSkillCategories() {
    const categories = await prisma.skills.findMany({
      where: { active: 1 },
      distinct: ['category'],
      select: {
        category: true
      },
      orderBy: {
        category: 'asc'
      }
    });

    return categories.map(cat => cat.category);
  }

  /**
   * Toggle skill active status
   */
  async toggleSkillStatus(skillId: number) {
    const skill = await prisma.skills.findUnique({
      where: { id: skillId }
    });

    if (!skill) {
      throw new ApiError(StatusCodes.NOT_FOUND, 'Skill not found');
    }

    return await prisma.skills.update({
      where: { id: skillId },
      data: {
        active: skill.active === 1 ? 0 : 1
      },
      include: {
        _count: {
          select: {
            user_skills: true,
            project_skills: true
          }
        }
      }
    });
  }

  /**
   * Get popular skills (most used by users)
   */
  async getPopularSkills(limit: number = 10) {
    const skills = await prisma.skills.findMany({
      where: { active: 1 },
      include: {
        _count: {
          select: {
            user_skills: true
          }
        }
      },
      orderBy: {
        user_skills: {
          _count: 'desc'
        }
      },
      take: limit
    });

    return skills.map(skill => ({
      ...skill,
      user_count: skill._count.user_skills
    }));
  }

  /**
   * Search skills by name or description
   */
  async searchSkills(query: string, limit: number = 20) {
    return await prisma.skills.findMany({
      where: {
        active: 1,
        OR: [
          { name: { contains: query } },
          { description: { contains: query } },
          { category: { contains: query } }
        ]
      },
      include: {
        _count: {
          select: {
            user_skills: true,
            project_skills: true
          }
        }
      },
      orderBy: {
        name: 'asc'
      },
      take: limit
    });
  }

  /**
   * Case-insensitive search helper (if needed for MySQL)
   */
  async searchSkillsCaseInsensitive(query: string, limit: number = 20) {


    const skills = await prisma.$queryRaw`
      SELECT s.*, 
        COUNT(us.user_id) as user_count,
        COUNT(ps.project_id) as project_count
      FROM skills s
      LEFT JOIN user_skills us ON s.id = us.skill_id
      LEFT JOIN project_skills ps ON s.id = ps.skill_id
      WHERE s.active = 1 
        AND (LOWER(s.name) LIKE LOWER(${`%${query}%`}) 
          OR LOWER(s.description) LIKE LOWER(${`%${query}%`}) 
          OR LOWER(s.category) LIKE LOWER(${`%${query}%`}))
      GROUP BY s.id
      ORDER BY s.name ASC
      LIMIT ${limit}
    `;

    return skills;
  }
}

export const skillsService = new SkillsService();
