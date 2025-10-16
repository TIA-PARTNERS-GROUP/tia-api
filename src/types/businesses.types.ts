import { z } from 'zod';

export const createBusinessSchema = z.object({
  operator_user_id: z.number()
    .int()
    .positive('Operator user ID must be a positive integer'),
  name: z.string()
    .min(1, 'Business name is required')
    .max(100, 'Business name must not exceed 100 characters')
    .trim(),
  tagline: z.string()
    .max(100, 'Tagline must not exceed 100 characters')
    .optional()
    .default(''),
  website: z.string()
    .max(255, 'Website must not exceed 255 characters')
    .url('Website must be a valid URL')
    .optional()
    .default(''),
  contact_name: z.string()
    .max(60, 'Contact name must not exceed 60 characters')
    .optional()
    .default(''),
  contact_phone_no: z.string()
    .max(20, 'Contact phone must not exceed 20 characters')
    .optional()
    .default(''),
  contact_email: z.string()
    .max(254, 'Contact email must not exceed 254 characters')
    .email('Contact email must be a valid email')
    .optional()
    .default(''),
  description: z.string()
    .max(5000, 'Description must not exceed 5000 characters')
    .optional()
    .default(''),
  address: z.string()
    .max(100, 'Address must not exceed 100 characters')
    .optional()
    .default(''),
  city: z.string()
    .max(60, 'City must not exceed 60 characters')
    .optional()
    .default(''),
  state: z.string()
    .max(60, 'State must not exceed 60 characters')
    .optional()
    .default(''),
  country: z.string()
    .max(60, 'Country must not exceed 60 characters')
    .optional()
    .default(''),
  postal_code: z.string()
    .max(20, 'Postal code must not exceed 20 characters')
    .optional()
    .default(''),
  value: z.number()
    .positive('Business value must be positive')
    .optional()
    .nullable(),
  business_type: z.enum(['Consulting', 'Retail', 'Technology', 'Manufacturing', 'Services', 'Other']),
  business_category: z.enum(['B2B', 'B2C', 'Non_Profit', 'Government', 'Mixed']),
  business_phase: z.enum(['Startup', 'Growth', 'Mature', 'Exit']),
  active: z.number()
    .int()
    .min(0)
    .max(1)
    .default(1)
});

export const updateBusinessSchema = z.object({
  name: z.string()
    .min(1, 'Business name is required')
    .max(100, 'Business name must not exceed 100 characters')
    .trim()
    .optional(),
  tagline: z.string()
    .max(100, 'Tagline must not exceed 100 characters')
    .optional()
    .nullable(),
  website: z.string()
    .max(255, 'Website must not exceed 255 characters')
    .url('Website must be a valid URL')
    .optional()
    .nullable(),
  contact_name: z.string()
    .max(60, 'Contact name must not exceed 60 characters')
    .optional()
    .nullable(),
  contact_phone_no: z.string()
    .max(20, 'Contact phone must not exceed 20 characters')
    .optional()
    .nullable(),
  contact_email: z.string()
    .max(254, 'Contact email must not exceed 254 characters')
    .email('Contact email must be a valid email')
    .optional()
    .nullable(),
  description: z.string()
    .max(5000, 'Description must not exceed 5000 characters')
    .optional()
    .nullable(),
  address: z.string()
    .max(100, 'Address must not exceed 100 characters')
    .optional()
    .nullable(),
  city: z.string()
    .max(60, 'City must not exceed 60 characters')
    .optional()
    .nullable(),
  state: z.string()
    .max(60, 'State must not exceed 60 characters')
    .optional()
    .nullable(),
  country: z.string()
    .max(60, 'Country must not exceed 60 characters')
    .optional()
    .nullable(),
  postal_code: z.string()
    .max(20, 'Postal code must not exceed 20 characters')
    .optional()
    .nullable(),
  value: z.number()
    .positive('Business value must be positive')
    .optional()
    .nullable(),
  business_type: z.enum(['Consulting', 'Retail', 'Technology', 'Manufacturing', 'Services', 'Other'])
    .optional(),
  business_category: z.enum(['B2B', 'B2C', 'Non_Profit', 'Government', 'Mixed'])
    .optional(),
  business_phase: z.enum(['Startup', 'Growth', 'Mature', 'Exit'])
    .optional(),
  active: z.number()
    .int()
    .min(0)
    .max(1)
    .optional()
}).refine(data => Object.keys(data).length > 0, {
  message: 'At least one field must be provided for update'
});

export const businessIdSchema = z.object({
  businessId: z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
    message: 'Business ID must be a positive integer'
  })
});

export const businessFilterSchema = z.object({
  business_type: z.string().optional(),
  business_category: z.string().optional(),
  business_phase: z.string().optional(),
  active: z.string().transform(val => val === 'true').optional(),
  search: z.string().optional(),
  operator_user_id: z.string().transform(val => parseInt(val, 10)).refine(val => !isNaN(val) && val > 0, {
    message: 'Operator user ID must be a positive integer'
  }).optional()
});

export type CreateBusinessInput = z.infer<typeof createBusinessSchema>;
export type UpdateBusinessInput = z.infer<typeof updateBusinessSchema>;
export type BusinessIdParams = z.infer<typeof businessIdSchema>;
export type BusinessFilterInput = z.infer<typeof businessFilterSchema>;
