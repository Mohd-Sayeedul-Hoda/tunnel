export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
  updated_at: string; // ISO date string
}

export interface APIKey {
  id: number;
  name: string;
  prefix: string;
  api_key_token?: string;
  user_id: number;
  expire_at: string; // ISO date string
  created_at: string; // ISO date string
  permission?: string[]; // Optional since it has omitempty
}

export type OtpType = 'email-verification' | 'forget-password';

export interface OtpVerification {
  email: string;
  otp: string;
  type: OtpType;
  expires_at: string; // ISO date string
  resend_count: number;
  created_at: string; // ISO date string
  updated_at: string; // ISO date string
}
