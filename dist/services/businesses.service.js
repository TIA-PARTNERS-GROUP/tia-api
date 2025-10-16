"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.businessesService = exports.BusinessesService = void 0;
const prisma_1 = require("../lib/prisma");
const ApiError_1 = require("../errors/ApiError");
const http_status_codes_1 = require("http-status-codes");
class BusinessesService {
    async getBusinesses(filters) {
        const where = {};
        if (filters?.business_type) {
            where.business_type = filters.business_type;
        }
        if (filters?.business_category) {
            where.business_category = filters.business_category;
        }
        if (filters?.business_phase) {
            where.business_phase = filters.business_phase;
        }
        if (filters?.active !== undefined) {
            where.active = filters.active ? 1 : 0;
        }
        if (filters?.operator_user_id) {
            where.operator_user_id = filters.operator_user_id;
        }
        if (filters?.search) {
            where.OR = [
                { name: { contains: filters.search } },
                { tagline: { contains: filters.search } },
                { description: { contains: filters.search } }
            ];
        }
        return await prisma_1.prisma.businesses.findMany({
            where,
            include: {
                users: {
                    select: {
                        id: true,
                        first_name: true,
                        last_name: true,
                        login_email: true
                    }
                },
                _count: {
                    select: {
                        business_connections_business_connections_initiating_business_idTobusinesses: true,
                        business_connections_business_connections_receiving_business_idTobusinesses: true,
                        projects: true,
                        publications: true,
                        business_tags: true
                    }
                }
            },
            orderBy: {
                name: 'asc'
            }
        });
    }
    async getBusinessById(businessId) {
        const business = await prisma_1.prisma.businesses.findUnique({
            where: { id: businessId },
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
                business_connections_business_connections_initiating_business_idTobusinesses: {
                    include: {
                        businesses_business_connections_receiving_business_idTobusinesses: {
                            select: {
                                id: true,
                                name: true,
                                business_type: true
                            }
                        }
                    },
                    take: 10
                },
                business_connections_business_connections_receiving_business_idTobusinesses: {
                    include: {
                        businesses_business_connections_initiating_business_idTobusinesses: {
                            select: {
                                id: true,
                                name: true,
                                business_type: true
                            }
                        }
                    },
                    take: 10
                },
                projects: {
                    select: {
                        id: true,
                        name: true,
                        project_status: true
                    },
                    take: 10
                },
                publications: {
                    select: {
                        id: true,
                        title: true,
                        publication_type: true,
                        published: true
                    },
                    take: 10
                },
                business_tags: {
                    select: {
                        id: true,
                        tag_type: true,
                        description: true
                    }
                }
            }
        });
        if (!business) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'Business not found');
        }
        return business;
    }
    async createBusiness(data) {
        const validatedData = data;
        const user = await prisma_1.prisma.users.findUnique({
            where: { id: validatedData.operator_user_id }
        });
        if (!user) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'Operator user not found');
        }
        return await prisma_1.prisma.businesses.create({
            data: {
                operator_user_id: validatedData.operator_user_id,
                name: validatedData.name,
                tagline: validatedData.tagline,
                website: validatedData.website,
                contact_name: validatedData.contact_name,
                contact_phone_no: validatedData.contact_phone_no,
                contact_email: validatedData.contact_email,
                description: validatedData.description,
                address: validatedData.address,
                city: validatedData.city,
                state: validatedData.state,
                country: validatedData.country,
                postal_code: validatedData.postal_code,
                value: validatedData.value,
                business_type: validatedData.business_type,
                business_category: validatedData.business_category,
                business_phase: validatedData.business_phase,
                active: validatedData.active
            },
            include: {
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
    async updateBusiness(businessId, data) {
        const validatedData = data;
        const existingBusiness = await prisma_1.prisma.businesses.findUnique({
            where: { id: businessId }
        });
        if (!existingBusiness) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'Business not found');
        }
        const updateData = {};
        if (validatedData.name !== undefined)
            updateData.name = validatedData.name;
        if (validatedData.tagline !== undefined)
            updateData.tagline = validatedData.tagline;
        if (validatedData.website !== undefined)
            updateData.website = validatedData.website;
        if (validatedData.contact_name !== undefined)
            updateData.contact_name = validatedData.contact_name;
        if (validatedData.contact_phone_no !== undefined)
            updateData.contact_phone_no = validatedData.contact_phone_no;
        if (validatedData.contact_email !== undefined)
            updateData.contact_email = validatedData.contact_email;
        if (validatedData.description !== undefined)
            updateData.description = validatedData.description;
        if (validatedData.address !== undefined)
            updateData.address = validatedData.address;
        if (validatedData.city !== undefined)
            updateData.city = validatedData.city;
        if (validatedData.state !== undefined)
            updateData.state = validatedData.state;
        if (validatedData.country !== undefined)
            updateData.country = validatedData.country;
        if (validatedData.postal_code !== undefined)
            updateData.postal_code = validatedData.postal_code;
        if (validatedData.value !== undefined)
            updateData.value = validatedData.value;
        if (validatedData.business_type !== undefined)
            updateData.business_type = validatedData.business_type;
        if (validatedData.business_category !== undefined)
            updateData.business_category = validatedData.business_category;
        if (validatedData.business_phase !== undefined)
            updateData.business_phase = validatedData.business_phase;
        if (validatedData.active !== undefined)
            updateData.active = validatedData.active;
        updateData.updated_at = new Date();
        return await prisma_1.prisma.businesses.update({
            where: { id: businessId },
            data: updateData,
            include: {
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
    async deleteBusiness(businessId) {
        const existingBusiness = await prisma_1.prisma.businesses.findUnique({
            where: { id: businessId }
        });
        if (!existingBusiness) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'Business not found');
        }
        const connectionsCount = await prisma_1.prisma.business_connections.count({
            where: {
                OR: [
                    { initiating_business_id: businessId },
                    { receiving_business_id: businessId }
                ]
            }
        });
        const projectsCount = await prisma_1.prisma.projects.count({
            where: { business_id: businessId }
        });
        const publicationsCount = await prisma_1.prisma.publications.count({
            where: { business_id: businessId }
        });
        if (connectionsCount > 0 || projectsCount > 0 || publicationsCount > 0) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.CONFLICT, `Cannot delete business. It has ${connectionsCount} connections, ${projectsCount} projects, and ${publicationsCount} publications.`);
        }
        await prisma_1.prisma.businesses.delete({
            where: { id: businessId }
        });
        return { message: 'Business deleted successfully' };
    }
    async getUserBusinesses(userId) {
        const user = await prisma_1.prisma.users.findUnique({
            where: { id: userId }
        });
        if (!user) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'User not found');
        }
        return await prisma_1.prisma.businesses.findMany({
            where: { operator_user_id: userId },
            include: {
                _count: {
                    select: {
                        projects: true,
                        publications: true,
                        business_tags: true
                    }
                }
            },
            orderBy: {
                name: 'asc'
            }
        });
    }
    async toggleBusinessStatus(businessId) {
        const business = await prisma_1.prisma.businesses.findUnique({
            where: { id: businessId }
        });
        if (!business) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'Business not found');
        }
        return await prisma_1.prisma.businesses.update({
            where: { id: businessId },
            data: {
                active: business.active === 1 ? 0 : 1,
                updated_at: new Date()
            },
            include: {
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
    async getBusinessStats(businessId) {
        const business = await prisma_1.prisma.businesses.findUnique({
            where: { id: businessId }
        });
        if (!business) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, 'Business not found');
        }
        const [outgoingConnections, incomingConnections, projects, publications, tags] = await Promise.all([
            prisma_1.prisma.business_connections.count({
                where: { initiating_business_id: businessId }
            }),
            prisma_1.prisma.business_connections.count({
                where: { receiving_business_id: businessId }
            }),
            prisma_1.prisma.projects.count({
                where: { business_id: businessId }
            }),
            prisma_1.prisma.publications.count({
                where: { business_id: businessId }
            }),
            prisma_1.prisma.business_tags.count({
                where: { business_id: businessId }
            })
        ]);
        return {
            business,
            stats: {
                outgoing_connections: outgoingConnections,
                incoming_connections: incomingConnections,
                projects: projects,
                publications: publications,
                tags: tags
            }
        };
    }
}
exports.BusinessesService = BusinessesService;
exports.businessesService = new BusinessesService();
