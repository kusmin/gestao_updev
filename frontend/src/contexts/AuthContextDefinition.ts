import { createContext } from 'react';
import { type LoginRequest, type SignupRequest } from '../lib/apiClient';

export type AuthState = {
  tenantId: string | null;
  userId: string | null;
  tokens: AuthTokens | null;
};

export type AuthTokens = {
  accessToken: string;
  refreshToken: string;
  expiresAt: number; // Timestamp
};

export type AuthContextValue = {
  isAuthenticated: boolean;
  initializing: boolean;
  tenantId: string | null;
  userId: string | null;
  accessToken: string | null;
  login: (credentials: LoginRequest) => Promise<void>;
  signup: (input: SignupRequest) => Promise<void>;
  logout: () => void;
};

export const AuthContext = createContext<AuthContextValue | undefined>(undefined);
