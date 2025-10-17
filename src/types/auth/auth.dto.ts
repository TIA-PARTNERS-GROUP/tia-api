export interface LoginResponse {
  user: {
    id: number;
    first_name: string;
    last_name: string | null;
    login_email: string;
    contact_email: string | null;
    email_verified: boolean;
    active: boolean;
    created_at: Date;
  };
  token: string;
  session_id: number;
  expires_at: Date;
  token_type: string;
}

export interface TokenValidationResponse {
  valid: boolean;
  user?: {
    id: number;
    first_name: string;
    last_name: string | null;
    login_email: string;
    email_verified: boolean;
    active: boolean;
  };
  session?: {
    id: number;
    ip_address: string | null;
    created_at: Date;
    expires_at: Date;
  };
}

export interface SessionInfo {
  id: number;
  ip_address: string | null;
  user_agent: string | null;
  created_at: Date;
  expires_at: Date;
  is_current?: boolean;
}

export interface LogoutResponse {
  message: string;
  sessions_ended: number;
}
