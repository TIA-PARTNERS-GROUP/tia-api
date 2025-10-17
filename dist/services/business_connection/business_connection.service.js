"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BusinessConnectionService = void 0;
const prisma_1 = require("../../lib/prisma");
const ApiError_1 = require("../../errors/ApiError");
const http_status_codes_1 = require("http-status-codes");
const mapToBusinessConnectionResponse = (prismaConnection) => {
    const { businesses_business_connections_initiating_business_idTobusinesses: initiating_business, businesses_business_connections_receiving_business_idTobusinesses: receiving_business, users: initiated_by_user, ...connectionData } = prismaConnection;
    const initiatingBusinessSummary = {
        id: initiating_business.id,
        name: initiating_business.name,
        business_type: initiating_business.business_type,
    };
    const receivingBusinessSummary = {
        id: receiving_business.id,
        name: receiving_business.name,
        business_type: receiving_business.business_type,
    };
    const initiatedByUserSummary = {
        id: initiated_by_user.id,
        first_name: initiated_by_user.first_name,
        last_name: initiated_by_user.last_name,
        login_email: initiated_by_user.login_email,
    };
    return {
        ...connectionData,
        initiating_business: initiatingBusinessSummary,
        receiving_business: receivingBusinessSummary,
        initiated_by_user: initiatedByUserSummary,
    };
};
class BusinessConnectionService {
    /**
     * Initiates a new connection request between two businesses.
     * Status is set to 'pending' by default.
     * * @param {InitiateConnectionInput} data The details of the connection to be established.
     * @returns {Promise<BusinessConnectionResponse>} The newly created connection record.
     * @throws {ApiError} 409 Conflict if the connection already exists.
     * @throws {ApiError} 404 Not Found if one of the business IDs or user ID is invalid.
     */
    async initiateConnection(data) {
        try {
            const connection = await prisma_1.prisma.business_connections.create({
                data: {
                    ...data,
                    status: "pending",
                },
                include: {
                    businesses_business_connections_initiating_business_idTobusinesses: true,
                    businesses_business_connections_receiving_business_idTobusinesses: true,
                    users: true,
                },
            });
            return mapToBusinessConnectionResponse(connection);
        }
        catch (error) {
            if (error.code === "P2002") {
                throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.CONFLICT, "A connection of this type already exists between these two businesses.");
            }
            if (error.code === "P2003") {
                throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, "One or both business IDs or the initiating user ID is invalid.");
            }
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR, "Failed to initiate business connection.");
        }
    }
    /**
     * Retrieves a single connection record by its unique ID.
     * * @param {number} connectionId The ID of the connection.
     * @returns {Promise<BusinessConnectionResponse | null>} The connection details.
     */
    async getConnectionById(connectionId) {
        const connection = await prisma_1.prisma.business_connections.findUnique({
            where: { id: connectionId },
            include: {
                businesses_business_connections_initiating_business_idTobusinesses: true,
                businesses_business_connections_receiving_business_idTobusinesses: true,
                users: true,
            },
        });
        if (!connection)
            return null;
        return mapToBusinessConnectionResponse(connection);
    }
    /**
     * Updates the status of an existing connection (e.g., from 'pending' to 'active').
     * * @param {number} connectionId The ID of the connection to update.
     * @param {UpdateConnectionStatusInput} data The new status and optional notes.
     * @returns {Promise<BusinessConnectionResponse>} The updated connection record.
     * @throws {ApiError} 404 Not Found if the connection ID is invalid.
     */
    async updateConnectionStatus(connectionId, data) {
        try {
            const updatedConnection = await prisma_1.prisma.business_connections.update({
                where: { id: connectionId },
                data: data,
                include: {
                    businesses_business_connections_initiating_business_idTobusinesses: true,
                    businesses_business_connections_receiving_business_idTobusinesses: true,
                    users: true,
                },
            });
            return mapToBusinessConnectionResponse(updatedConnection);
        }
        catch (error) {
            if (error.code === "P2025") {
                throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.NOT_FOUND, "Business connection not found.");
            }
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR, "Failed to update connection status.");
        }
    }
    /**
     * Retrieves all connections where the given business is either the initiator or receiver.
     * * @param {number} businessId The ID of the business whose connections are being fetched.
     * @param {'pending' | 'active' | 'rejected' | 'inactive'} [status] Optional filter by connection status.
     * @returns {Promise<ConnectionSummaryResponse[]>} A list of connection summaries.
     */
    async getConnectionsForBusiness(businessId, status) {
        try {
            const connections = await prisma_1.prisma.business_connections.findMany({
                where: {
                    OR: [
                        { initiating_business_id: businessId },
                        { receiving_business_id: businessId },
                    ],
                    ...(status && { status: status }),
                },
                select: {
                    id: true,
                    initiating_business_id: true,
                    receiving_business_id: true,
                    connection_type: true,
                    status: true,
                    created_at: true,
                },
                orderBy: { created_at: "desc" },
            });
            return connections;
        }
        catch (error) {
            throw new ApiError_1.ApiError(http_status_codes_1.StatusCodes.INTERNAL_SERVER_ERROR, "Failed to retrieve business connections.");
        }
    }
}
exports.BusinessConnectionService = BusinessConnectionService;
