"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.businessFilterSchema = exports.businessIdSchema = exports.updateBusinessSchema = exports.createBusinessSchema = void 0;
const zod_1 = require("zod");
exports.createBusinessSchema = zod_1.z.object({
    operator_user_id: zod_1.z.number()
        .int()
        .positive('Operator user ID must be a positive integer'),
    name: zod_1.z.string()
        .min(1, 'Business name is required')
        .max(100, 'Business name must not exceed 100 characters')
        .trim(),
    tagline: zod_1.z.string()
        .max(100, 'Tagline must not exceed 100 characters')
        .optional()
        .default(''),
    website: zod_1.z.string()
        .max(255, 'Website must not exceed 255 characters')
        .url('Website must be a valid URL')
        .optional()
        .default(''),
    contact_name: zod_1.z.string()
        .max(60, 'Contact name must not exceed 60 characters')
        .optional()
        .default(''),
    contact_phone_no: zod_1.z.string()
        .max(20, 'Contact phone must not exceed 20 characters')
        .optional()
        .default(''),
    contact_email: zod_1.z.string()
        .max(254, 'Contact email must not exceed 254 characters')
        .email('Contact email must be a valid email')
        .optional()
        .default(''),
    description: zod_1.z.string()
        .max(5000, 'Description must not exceed 5000 characters')
        .optional()
        .default(''),
    address: zod_1.z.string()
        .max(100, 'Address must not exceed 100 characters')
        .optional()
        .default(''),
    city: zod_1.z.string()
        .max(60, 'City must not exceed 60 characters')
        .optional()
        .default(''),
    state: zod_1.z.string()
        .max(60, 'State must not exceed 60 characters')
        .optional()
        .default(''),
    country: zod_1.z.string()
        .max(60, 'Country must not exceed 60 characters')
        .optional()
        .default(''),
    postal_code: zod_1.z.string()
        .max(20, 'Postal code must not exceed 20 characters')
        .optional()
        .default(''),
    value: zod_1.z.number()
        .positive('Business value must be positive')
        .optional()
        .nullable(),
    business_type: zod_1.z.enum(['Consulting', 'Retail', 'Technology', 'Manufacturing', 'Services', 'Other']),
    business_category: zod_1.z.enum(['B2B', 'B2C', 'Non_Profit', 'Government', 'Mixed']),
    business_phase: zod_1.z.enum(['Startup', 'Growth', 'Mature', 'Exit']),
    active: zod_1.z.number()
        .int()
        .min(0)
        .max(1)
        .default(1)
});
exports.updateBusinessSchema = zod_1.z.object({
    name: zod_1.z.string()
        .min(1, 'Business name is required')
        .max(100, 'Business name must not exceed 100 characters')
        .trim()
        .optional(),
    tagline: zod_1.z.string()
        .max(100, 'Tagline must not exceed 100 characters')
        .optional()
        .nullable(),
    website: zod_1.z.string()
        .max(255, 'Website must not exceed 255 characters')
        .url('Website must be a valid URL')
        .optional()
        .nullable(),
    contact_name: zod_1.z.string()
        .max(60, 'Contact name must not exceed 60 characters')
        .optional()
        .nullable(),
    contact_phone_no: zod_1.z.string()
        .max(20, 'Contact phone must not exceed 20 characters')
        .optional()
        .nullable(),
    contact_email: zod_1.z.string()
        .max(254, 'Contact email must not exceed 254 characters')
        .email('Contact email must be a valid email')
        .optional()
        .nullable(),
    description: zod_1.z.string()
        .max(5000, 'Description must not exceed 5000 characters')
        .optional()
        .nullable(),
    address: zod_1.z.string()
        .max(100, 'Address must not exceed 100 characters')
        .optional()
        .nullable(),
    city: zod_1.z.string()
        .max(60, 'City must not exceed 60 characters')
        .optional()
        .nullable(),
    state: zod_1.z.string()
        .max(60, 'State must not exceed 60 characters')
        .optional()
        .nullable(),
    country: zod_1.z.string()
        .max(60, 'Country must not exceed 60 characters')
        .optional()
        .nullable(),
    postal_code: zod_1.z.string()
        .max(20, 'Postal code must not exceed 20 characters')
        .optional()
        .nullable(),
    value: zod_1.z.number()
        .positive('Business value must be positive')
        .optional()
        .nullable(),
    business_type: zod_1.z.enum(['Consulting', 'Retail', 'Technology', 'Manufacturing', 'Services', 'Other'])
        .optional(),
    business_category: zod_1.z.enum(['B2B', 'B2C', 'Non_Profit', 'Government', 'Mixed'])
        .optional(),
    business_phase: zod_1.z.enum(['Startup', 'Growth', 'Mature', 'Exit'])
        .optional(),
    active: zod_1.z.number()
        .int()
        .min(0)
        .max(1)
        .optional()
}).refine(data => Object.keys(data).length > 0, {
    message: 'At least one field must be provided for update'
});
exports.businessIdSchema = zod_1.z.object({
    businessId: zod_1.z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
        message: 'Business ID must be a positive integer'
    })
});
exports.businessFilterSchema = zod_1.z.object({
    business_type: zod_1.z.string().optional(),
    business_category: zod_1.z.string().optional(),
    business_phase: zod_1.z.string().optional(),
    active: zod_1.z.string().transform(val => val === 'true').optional(),
    search: zod_1.z.string().optional(),
    operator_user_id: zod_1.z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
        message: 'Operator user ID must be a positive integer'
    }).optional()
});
