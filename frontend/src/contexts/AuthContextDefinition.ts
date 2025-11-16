import { createContext } from 'react';
import { type LoginRequest, type SignupRequest } from '../lib/apiClient';

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
