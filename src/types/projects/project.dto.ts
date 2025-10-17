/**
 * Summary representation of a user, used for nested relationships
 * to provide context without exposing full user details.
 */
interface UserSummaryResponse {
  /** The unique identifier of the user. */
  id: number;
  /** The user's first name. */
  first_name: string;
  /** The user's last name. */
  last_name: string | null;
  /** The user's primary login email address. */
  login_email: string;
}

/**
 * Summary representation of a business, used for nested relationships.
 */
interface BusinessSummaryResponse {
  /** The unique identifier of the business. */
  id: number;
  /** The official name of the business. */
  name: string;
  /** The primary type of the business (e.g., Consulting, Retail). */
  business_type: string;
}

/**
 * Detailed representation of a user's membership within a specific project.
 * Used for responses related to project member management.
 */
export interface ProjectMemberResponse {
  /** The unique identifier of the project the user belongs to. */
  project_id: number;
  /** The unique identifier of the user who is the member. */
  user_id: number;
  /** The role of the user within the project. */
  role: "manager" | "contributor" | "reviewer";
  /** The timestamp when the user joined the project. */
  joined_at: Date;
  /** Summary details of the project member. */
  user: UserSummaryResponse;
}

/**
 * Complete and detailed representation of a single project, including
 * its manager, associated business, and list of members.
 */
export interface ProjectResponse {
  /** The unique identifier for the project. */
  id: number;
  /** The ID of the user responsible for managing the project. */
  managed_by_user_id: number;
  /** The ID of the associated business, or null if it's a personal project. */
  business_id: number | null;
  /** The official name of the project. */
  name: string;
  /** A detailed description of the project's goals and scope. */
  description: string | null;
  /** The current phase or status of the project. */
  project_status: "planning" | "active" | "on_hold" | "completed" | "cancelled";
  /** The planned start date of the project. */
  start_date: Date | null;
  /** The target or planned end date of the project. */
  target_end_date: Date | null;
  /** The actual date the project was completed (null if not completed). */
  actual_end_date: Date | null;
  /** The timestamp when the project record was created. */
  created_at: Date;
  /** The timestamp of the last update to the project record. */
  updated_at: Date;

  /** Summary details of the project manager. */
  manager: UserSummaryResponse;
  /** Summary details of the associated business, if applicable. */
  business?: BusinessSummaryResponse | null;
  /** A list of all members currently assigned to the project. */
  members: ProjectMemberResponse[];
}

/**
 * A light-weight representation of a project, suitable for lists or summaries.
 */
export interface ProjectSummaryResponse {
  /** The unique identifier for the project. */
  id: number;
  /** The official name of the project. */
  name: string;
  /** The current phase or status of the project. */
  project_status: string;
  /** The planned start date of the project. */
  start_date: Date | null;
}
