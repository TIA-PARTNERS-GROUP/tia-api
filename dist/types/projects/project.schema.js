"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.updateMemberRoleSchema = exports.addMemberSchema = exports.projectIdParamsSchema = exports.updateProjectSchema = exports.createProjectSchema = void 0;
const zod_1 = require("zod");
/**
 * Defines the possible status values for a project.
 */
const ProjectStatusEnum = zod_1.z.enum([
    "planning",
    "active",
    "on_hold",
    "completed",
    "cancelled",
]);
/**
 * Defines the possible roles for a user within a project.
 */
const ProjectMemberRoleEnum = zod_1.z.enum(["manager", "contributor", "reviewer"]);
/**
 * Schema for validating data required to create a new project.
 * All fields are mandatory unless marked optional/nullable.
 */
exports.createProjectSchema = zod_1.z.object({
    /** The unique ID of the user managing the project. Must be a positive integer. */
    managed_by_user_id: zod_1.z
        .number()
        .int()
        .positive("Manager ID must be a positive integer"),
    /** The unique ID of the associated business, or null if it's a personal project. */
    business_id: zod_1.z
        .number()
        .int()
        .positive("Business ID must be a positive integer")
        .optional()
        .nullable(),
    /** The name of the project. Required, trimmed, and limited to 100 characters. */
    name: zod_1.z
        .string()
        .trim()
        .min(1, "Project name is required")
        .max(100, "Name cannot exceed 100 characters"),
    /** A detailed description of the project. Optional, nullable, and limited to 1000 characters. */
    description: zod_1.z
        .string()
        .max(1000, "Description cannot exceed 1000 characters")
        .optional()
        .nullable(),
    /** The initial status of the project. Defaults to 'planning'. */
    project_status: ProjectStatusEnum.default("planning"),
    /** The planned start date of the project. Optional and nullable. */
    start_date: zod_1.z.date().optional().nullable(),
    /** The target or planned end date of the project. Optional and nullable. */
    target_end_date: zod_1.z.date().optional().nullable(),
});
/**
 * Schema for validating data used to partially update an existing project.
 * All fields are optional, but at least one field must be present.
 */
exports.updateProjectSchema = exports.createProjectSchema
    .partial()
    .extend({
    /** The actual date the project was completed. Optional and nullable. */
    actual_end_date: zod_1.z.date().optional().nullable(),
})
    .refine((data) => Object.keys(data).length > 0, {
    message: "At least one field must be provided for update",
});
/**
 * Schema for validating and transforming the 'projectId' path parameter.
 */
exports.projectIdParamsSchema = zod_1.z.object({
    /** The project ID retrieved from the URL path, transformed into a positive integer. */
    projectId: zod_1.z
        .string()
        .transform((val) => parseInt(val, 10))
        .refine((val) => !isNaN(val) && val > 0, {
        message: "Project ID must be a positive integer",
    }),
});
/**
 * Schema for validating data required to add a new member to a project.
 */
exports.addMemberSchema = zod_1.z.object({
    /** The unique ID of the user to be added as a member. Must be a positive integer. */
    user_id: zod_1.z.number().int().positive("User ID must be a positive integer"),
    /** The role to assign to the new member. Defaults to 'contributor'. */
    role: ProjectMemberRoleEnum.default("contributor"),
});
/**
 * Schema for validating data required to update an existing member's role.
 */
exports.updateMemberRoleSchema = zod_1.z.object({
    /** The new role for the project member. */
    role: ProjectMemberRoleEnum,
});
