import type {
  CreateBusinessInput,
  UpdateBusinessInput,
} from "./businesses.validation.js";

export type { CreateBusinessInput, UpdateBusinessInput };

/**
 * Defines the shape of the business data returned by the API.
 * This is the primary Data Transfer Object (DTO) for a business profile.
 */
export interface BusinessResponse {
  /**
   * The unique identifier for the business.
   * @example 123
   */
  id: number;

  /**
   * The ID of the user who operates or owns the business.
   * @example 45
   */
  operator_user_id: number;

  /**
   * The official name of the business.
   * @example "Innovate Solutions Inc."
   */
  name: string;

  /**
   * A short, catchy phrase describing the business.
   * @example "Pioneering the Future of Tech"
   */
  tagline: string | null;

  /**
   * The official website URL of the business.
   * @example "https://innovatesolutions.com"
   */
  website: string | null;

  /**
   * The name of the primary contact person for the business.
   * @example "Jane Doe"
   */
  contact_name: string | null;

  /**
   * The primary contact phone number.
   * @example "+1-555-123-4567"
   */
  contact_phone_no: string | null;

  /**
   * The primary contact email address.
   * @example "contact@innovatesolutions.com"
   */
  contact_email: string | null;

  /**
   * A detailed description of the business, its mission, and services.
   */
  description: string | null;

  /**
   * The street address of the business.
   * @example "123 Innovation Drive"
   */
  address: string | null;

  /**
   * The city where the business is located.
   * @example "Metropolis"
   */
  city: string | null;

  /**
   * The state or province.
   * @example "California"
   */
  state: string | null;

  /**
   * The country where the business is located.
   * @example "USA"
   */
  country: string | null;

  /**
   * The postal or ZIP code.
   * @example "90210"
   */
  postal_code: string | null;

  /**
   * The estimated monetary value of the business.
   * @example 5000000
   */
  value: number | null;

  /**
   * The primary industry or type of the business.
   */
  business_type:
    | "Consulting"
    | "Retail"
    | "Technology"
    | "Manufacturing"
    | "Services"
    | "Other";

  /**
   * The target market category.
   */
  business_category: "B2B" | "B2C" | "Non_Profit" | "Government" | "Mixed";

  /**
   * The current developmental stage of the business.
   */
  business_phase: "Startup" | "Growth" | "Mature" | "Exit";

  /**
   * Indicates if the business profile is active.
   * @example true
   */
  active: boolean;

  /**
   * The timestamp when the business profile was created.
   */
  created_at: Date;

  /**
   * The timestamp when the business profile was last updated.
   */
  updated_at: Date;
}

/**
 * Defines the shape for filtering business queries.
 * Used internally by the service layer.
 */
export interface BusinessesFilter {
  business_type?: string;
  business_category?: string;
  business_phase?: string;
  active?: boolean;
  search?: string;
  operator_user_id?: number;
}

/**
 * Defines the shape of the statistics data returned for a business.
 */
export interface BusinessStatsResponse {
  /**
   * The total number of projects associated with the business.
   * @example 25
   */
  total_projects: number;

  /**
   * The number of active team members.
   * @example 15
   */
  team_member_count: number;

  /**
   * The overall completion rate of projects.
   * @example 0.85
   */
  project_completion_rate: number;
}
