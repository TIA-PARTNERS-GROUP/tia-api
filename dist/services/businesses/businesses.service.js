"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BusinessesService = void 0;
const prisma_1 = require("../../lib/prisma");
const ApiError_1 = require("../../errors/ApiError");
const http_status_codes_1 = require("http-status-codes");
/**
 * A reusable function to transform the raw business data from Prisma
 * into the clean BusinessResponse DTO that the API will return.
 * @param business - The raw business object from a Prisma query.
 * @returns A BusinessResponse DTO.
 */
const transformBusinessToDto = (business) => {
    return {
        ...business,
        value: business.value ? business.value.toNumber() : null,
        active: business.active === 1,
    };
};
class BusinessesService {
    /**
     * Retrieves a filtered and paginated list of businesses.
     */
    async getBusinesses(filters) {
        const where = {};
        if (filters?.business_type)
            where.business_type = filters.business_type;
        if (filters?.business_category)
            where.business_category = filters.business_category;
        if (filters?.business_phase)
            where.business_phase = filters.business_phase;
        if (filters?.active !== undefined)
            where.active = filters.active ? 1 : 0;
        if (filters?.operator_user_id)
            where.operator_user_id = filters.operator_user_id;
        if (filters?.search) {
            where.OR = [
                { name: { contains: filters.search } },
                { tagline: { contains: filters.search } },
                { description: { contains: filters.search } },
            ];
        }
        const businesses = await prisma_1.prisma.businesses.findMany({
            where,
            include: {
                users: {
                    select: {
                        id: true,
                        first_name: true,
                        last_name: true,
                        login_email: true,
                    },
                },
            },
            orderBy: { name: "asc" },
        });
        return businesses.map(transformBusinessToDto);
    }
    /**
     * Retrieves a single business by its unique ID.
     */
    async getBusinessById(businessId) {
        const business = await prisma_1.prisma.businesses.findUnique({
            where: { id: businessId },
            include: {
                users: true,
            },
        });
        if (!business) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, "Business not found");
        }
        return transformBusinessToDto(business);
    }
    /**
     * Creates a new business.
     */
    async createBusiness(data) {
        const user = await prisma_1.prisma.users.findUnique({
            where: { id: data.operator_user_id },
        });
        if (!user) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, "Operator user not found");
        }
        const newBusiness = await prisma_1.prisma.businesses.create({
            data: {
                ...data,
                value: data.value,
            },
        });
        return transformBusinessToDto(newBusiness);
    }
    /**
     * Updates an existing business.
     */
    async updateBusiness(businessId, data) {
        const existingBusiness = await prisma_1.prisma.businesses.findUnique({
            where: { id: businessId },
        });
        if (!existingBusiness) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, "Business not found");
        }
        const updatedBusiness = await prisma_1.prisma.businesses.update({
            where: { id: businessId },
            data: {
                ...data,
                value: data.value,
                updated_at: new Date(),
            },
        });
        return transformBusinessToDto(updatedBusiness);
    }
    /**
     * Deletes a business if it has no dependencies.
     */
    async deleteBusiness(businessId) {
        const business = await prisma_1.prisma.businesses.findUnique({
            where: { id: businessId },
            include: { _count: { select: { projects: true, publications: true } } },
        });
        if (!business) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, "Business not found");
        }
        if (business._count.projects > 0 || business._count.publications > 0) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.CONFLICT, `Cannot delete business. It is associated with ${business._count.projects} projects and ${business._count.publications} publications.`);
        }
        await prisma_1.prisma.businesses.delete({
            where: { id: businessId },
        });
    }
    /**
     * Retrieves all businesses operated by a specific user.
     */
    async getUserBusinesses(userId) {
        const user = await prisma_1.prisma.users.findUnique({ where: { id: userId } });
        if (!user)
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, "User not found");
        const businesses = await prisma_1.prisma.businesses.findMany({
            where: { operator_user_id: userId },
            orderBy: { name: "asc" },
        });
        return businesses.map(transformBusinessToDto);
    }
    /**
     * Toggles the active status of a business.
     */
    async toggleBusinessStatus(businessId) {
        const business = await prisma_1.prisma.businesses.findUnique({
            where: { id: businessId },
        });
        if (!business) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, "Business not found");
        }
        const updatedBusiness = await prisma_1.prisma.businesses.update({
            where: { id: businessId },
            data: {
                active: business.active === 1 ? 0 : 1,
                updated_at: new Date(),
            },
        });
        return transformBusinessToDto(updatedBusiness);
    }
    /**
     * Retrieves key statistics for a business.
     */
    async getBusinessStats(businessId) {
        const stats = await prisma_1.prisma.businesses.findUnique({
            where: { id: businessId },
            select: {
                _count: {
                    select: {
                        projects: true,
                        publications: true,
                        business_tags: true,
                    },
                },
            },
        });
        if (!stats) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, "Business not found");
        }
        return {
            total_projects: stats._count.projects,
            team_member_count: 0,
            project_completion_rate: 0,
        };
    }
}
exports.BusinessesService = BusinessesService;
