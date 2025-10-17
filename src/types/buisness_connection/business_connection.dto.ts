
interface BusinessSummaryResponse {
  id: number;
  name: string;
  business_type: string;
}

interface UserSummaryResponse {
  id: number;
  first_name: string;
  last_name: string | null;
  login_email: string;
}

/**
 * Detailed representation of a single business connection record.
 * Used for detailed responses after creation or retrieval by ID.
 */
export interface BusinessConnectionResponse {
  /** The unique identifier of the connection record. */
  id: number;
  /** The ID of the business that initiated the connection request. */
  initiating_business_id: number;
  /** The ID of the business that received the connection request. */
  receiving_business_id: number;
  /** The defined type of business relationship. */
  connection_type:
    | "Partnership"
    | "Supplier"
    | "Client"
    | "Referral"
    | "Collaboration";
  /** The current status of the connection (e.g., 'pending', 'active'). */
  status: "pending" | "active" | "rejected" | "inactive";
  /** The ID of the user who performed the initiation (for audit). */
  initiated_by_user_id: number;
  /** Notes or message associated with the connection request or status change. */
  notes: string | null;
  /** The timestamp when the connection record was created. */
  created_at: Date;
  /** The timestamp of the last update to the connection record. */
  updated_at: Date;

 
  /** Summary details of the initiating business. */
  initiating_business: BusinessSummaryResponse;
  /** Summary details of the receiving business. */
  receiving_business: BusinessSummaryResponse;
  /** Summary details of the user who initiated the connection. */
  initiated_by_user: UserSummaryResponse;
}

/**
 * Light-weight representation for listing or summary views of business connections.
 */
export interface ConnectionSummaryResponse {
  /** The unique identifier of the connection record. */
  id: number;
  /** The ID of the initiating business. */
  initiating_business_id: number;
  /** The ID of the receiving business. */
  receiving_business_id: number;
  /** The defined type of business relationship. */
  connection_type: string;
  /** The current status of the connection. */
  status: string;
  /** The timestamp when the connection record was created. */
  created_at: Date;
}
