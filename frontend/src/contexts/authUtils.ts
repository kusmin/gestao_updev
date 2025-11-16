import {
  type AuthTokensResponse,
  type SignupResponse,
} from '../lib/apiClient';
import { AUTH_STORAGE_KEY } from '../utils/constants';

type AuthTokens = {
  accessToken: string;
  refreshToken: string;
  expiresAt: number;
};

type AuthState = {
  tenantId: string | null;
  userId: string | null;
  tokens: AuthTokens | null;
};

export const DEFAULT_STATE: AuthState = {
  tenantId: null,
  userId: null,
  tokens: null,
};

export const safeAtob = (value: string) => {
  try {
    if (typeof atob === 'function') {
      return atob(value);
    }
    return Buffer.from(value, 'base64').toString('utf-8');
  } catch {
    return '';
  }
};

export const decodeTokenClaims = (token: string): { tenantId?: string; userId?: string } => {
  try {
    const parts = token.split('.');
    if (parts.length < 2) {
      return {};
    }
    const base64 = parts[1].replace(/-/g, '+').replace(/_/g, '/');
    const payload = JSON.parse(safeAtob(base64));
    return {
      tenantId: payload?.tenant_id ?? payload?.tenantId,
      userId: payload?.user_id ?? payload?.userId,
    };
  } catch {
    return {};
  }
};

export const mapTokens = (
  payload: AuthTokensResponse,
  fallbackTenant?: string | null,
  fallbackUser?: string | null,
): AuthState => {
  const claims = decodeTokenClaims(payload.access_token);
  const tenantId = payload.tenant_id ?? claims.tenantId ?? fallbackTenant ?? null;
  const userId = payload.user_id ?? claims.userId ?? fallbackUser ?? null;
  if (!tenantId) {
    throw new Error('Tenant ID ausente da resposta de autenticação.');
  }
  return {
    tenantId,
    userId,
    tokens: {
      accessToken: payload.access_token,
      refreshToken: payload.refresh_token,
      expiresAt: Date.now() + payload.expires_in * 1000,
    },
  };
};

export const mapSignupResult = (result: SignupResponse): AuthState => ({
  tenantId: result.tenant_id,
  userId: result.user_id,
  tokens: {
    accessToken: result.access_token,
    refreshToken: result.refresh_token,
    expiresAt: Date.now() + result.expires_in * 1000,
  },
});

export const loadStoredAuth = (): AuthState => {
  if (typeof window === 'undefined') {
    return DEFAULT_STATE;
  }
  const raw = window.localStorage.getItem(AUTH_STORAGE_KEY);
  if (!raw) {
    return DEFAULT_STATE;
  }
  try {
    const parsed = JSON.parse(raw) as AuthState;
    if (parsed.tokens && parsed.tokens.expiresAt <= Date.now()) {
      return DEFAULT_STATE;
    }
    return {
      tenantId: parsed.tenantId ?? null,
      userId: parsed.userId ?? null,
      tokens: parsed.tokens ?? null,
    };
  } catch {
    return DEFAULT_STATE;
  }
};

export const persistState = (state: AuthState) => {
  if (typeof window === 'undefined') {
    return;
  }
  if (!state.tokens) {
    window.localStorage.removeItem(AUTH_STORAGE_KEY);
    return;
  }
  window.localStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(state));
};
