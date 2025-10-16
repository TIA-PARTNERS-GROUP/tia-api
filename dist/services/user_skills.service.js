"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.userSkillsService = exports.UserSkillsService = void 0;
const prisma_1 = require("../lib/prisma");
const user_skills_validation_1 = require("../types/user_skills.validation");
const ApiError_1 = require("../errors/ApiError");
const http_status_codes_1 = require("http-status-codes");
class UserSkillsService {
    /**
     * Get all skills for a user
     */
    async getUserSkills(userId) {
        return await prisma_1.prisma.user_skills.findMany({
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
    async getUserSkillById(skillId, userId) {
        const userSkill = await prisma_1.prisma.user_skills.findUnique({
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
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'User skill not found');
        }
        return userSkill;
    }
    /**
     * Create a new user skill
     */
    async createUserSkill(data) {
        const validatedData = user_skills_validation_1.createUserSkillSchema.parse(data);
        const user = await prisma_1.prisma.users.findUnique({
            where: { id: validatedData.user_id }
        });
        if (!user) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'User not found');
        }
        const skill = await prisma_1.prisma.skills.findUnique({
            where: { id: validatedData.skill_id }
        });
        if (!skill) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'Skill not found');
        }
        const existingUserSkill = await prisma_1.prisma.user_skills.findUnique({
            where: {
                skill_id_user_id: {
                    skill_id: validatedData.skill_id,
                    user_id: validatedData.user_id
                }
            }
        });
        if (existingUserSkill) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.CONFLICT, 'User already has this skill');
        }
        return await prisma_1.prisma.user_skills.create({
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
    async updateUserSkill(skillId, userId, data) {
        const validatedData = user_skills_validation_1.updateUserSkillSchema.parse(data);
        const existingUserSkill = await prisma_1.prisma.user_skills.findUnique({
            where: {
                skill_id_user_id: {
                    skill_id: skillId,
                    user_id: userId
                }
            }
        });
        if (!existingUserSkill) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'User skill not found');
        }
        return await prisma_1.prisma.user_skills.update({
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
    async deleteUserSkill(skillId, userId) {
        const existingUserSkill = await prisma_1.prisma.user_skills.findUnique({
            where: {
                skill_id_user_id: {
                    skill_id: skillId,
                    user_id: userId
                }
            }
        });
        if (!existingUserSkill) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'User skill not found');
        }
        await prisma_1.prisma.user_skills.delete({
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
    async getUsersBySkill(skillId) {
        const skill = await prisma_1.prisma.skills.findUnique({
            where: { id: skillId }
        });
        if (!skill) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'Skill not found');
        }
        return await prisma_1.prisma.user_skills.findMany({
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
    async getUserSkillsByProficiency(userId, proficiencyLevel) {
        return await prisma_1.prisma.user_skills.findMany({
            where: {
                user_id: userId,
                proficiency_level: proficiencyLevel
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
exports.UserSkillsService = UserSkillsService;
exports.userSkillsService = new UserSkillsService();
