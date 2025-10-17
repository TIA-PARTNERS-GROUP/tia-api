export interface UserResponse {
  id: number;
  first_name: string;
  last_name: string | null;
  login_email: string;
  contact_email: string | null;
  contact_phone_no: string | null;
  adk_session_id: string | null;
  email_verified: boolean;
  active: boolean;
  created_at: Date;
  updated_at: Date;
}

export interface UserAuthResponse {
  id: number;
  first_name: string;
  last_name: string | null;
  login_email: string;
  contact_email: string | null;
  email_verified: boolean;
  active: boolean;
  created_at: Date;
}

export interface PasswordResetRequest {
  token: string;
  new_password: string;
}

export interface PasswordResetInitiateRequest {
  login_email: string;
}
