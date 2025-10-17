"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.updateConnectionStatusSchema = exports.connectionIdParamsSchema = exports.initiateConnectionSchema = void 0;
const zod_1 = require("zod");
const ConnectionTypeEnum = zod_1.z.enum([
    "Partnership",
    "Supplier",
    "Client",
    "Referral",
    "Collaboration",
]);
const ConnectionStatusEnum = zod_1.z.enum([
    "pending",
    "active",
    "rejected",
    "inactive",
]);
/**
 * Schema for validating data required to initiate a new business connection request.
 */
exports.initiateConnectionSchema = zod_1.z
    .object({
    /** The ID of the business sending the connection request. Must be a positive integer. */
    initiating_business_id: zod_1.z
        .number()
        .int()
        .positive("Initiating Business ID must be a positive integer"),
    /** The ID of the business receiving the connection request. Must be a positive integer. */
    receiving_business_id: zod_1.z
        .number()
        .int()
        .positive("Receiving Business ID must be a positive integer"),
    /** The type of connection being established. */
    connection_type: ConnectionTypeEnum,
    /** The ID of the user initiating the request. Used for logging and permissions. Must be a positive integer. */
    initiated_by_user_id: zod_1.z
        .number()
        .int()
        .positive("Initiated By User ID must be a positive integer"),
    /** Optional notes or message included with the connection request. Max 1000 characters. */
    notes: zod_1.z.string().max(1000).optional().nullable(),
})
    .refine((data) => data.initiating_business_id !== data.receiving_business_id, {
    message: "A business cannot initiate a connection with itself.",
    path: ["receiving_business_id"],
});
/**
 * Schema for validating and transforming the connection ID path parameter.
 */
exports.connectionIdParamsSchema = zod_1.z.object({
    /** The Connection ID retrieved from the URL path, transformed into a positive integer. */
    connectionId: zod_1.z
        .string()
        .transform((val) => parseInt(val, 10))
        .refine((val) => !isNaN(val) && val > 0, {
        message: "Connection ID must be a positive integer",
    }),
});
/**
 * Schema for validating data required to update the status of a connection.
 */
exports.updateConnectionStatusSchema = zod_1.z.object({
    /** The new status for the connection (e.g., 'active' or 'rejected'). */
    status: ConnectionStatusEnum,
    /** Optional: Notes explaining the status change (e.g., reason for rejection). Max 1000 characters. */
    notes: zod_1.z.string().max(1000).optional().nullable(),
});
