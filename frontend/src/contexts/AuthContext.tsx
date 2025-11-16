import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react';
import {
  login as loginRequest,
  signup as signupRequest,
  refreshTokens as refreshTokensRequest,
  type LoginRequest,
  type SignupRequest,
  type AuthTokensResponse,
  type SignupResponse,
} from '../lib/apiClient';

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

const DEFAULT_STATE: AuthState = {
  tenantId: null,
  userId: null,
  tokens: null,
};

type AuthContextValue = {
  isAuthenticated: boolean;
  initializing: boolean;
  tenantId: string | null;
  userId: string | null;
  accessToken: string | null;
  login: (credentials: LoginRequest) => Promise<void>;
  signup: (input: SignupRequest) => Promise<void>;
  logout: () => void;
};

export const AUTH_STORAGE_KEY = 'gestao-auth';

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

const safeAtob = (value: string) => {
  try {
    if (typeof atob === 'function') {
      return atob(value);
    }
    return Buffer.from(value, 'base64').toString('utf-8');
  } catch {
    return '';
  }
};

const decodeTokenClaims = (token: string): { tenantId?: string; userId?: string } => {
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

const mapTokens = (
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

const mapSignupResult = (result: SignupResponse): AuthState => ({
  tenantId: result.tenant_id,
  userId: result.user_id,
  tokens: {
    accessToken: result.access_token,
    refreshToken: result.refresh_token,
    expiresAt: Date.now() + result.expires_in * 1000,
  },
});

const loadStoredAuth = (): AuthState => {
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

const persistState = (state: AuthState) => {
  if (typeof window === 'undefined') {
    return;
  }
  if (!state.tokens) {
    window.localStorage.removeItem(AUTH_STORAGE_KEY);
    return;
  }
  window.localStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(state));
};

type AuthProviderProps = {
  children: React.ReactNode;
};

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [state, setState] = useState<AuthState>(() => loadStoredAuth());
  const [initializing, setInitializing] = useState(true);
  const refreshTimer = useRef<number | null>(null);

  const clearRefreshTimer = () => {
    if (refreshTimer.current) {
      window.clearTimeout(refreshTimer.current);
      refreshTimer.current = null;
    }
  };

  const scheduleRefresh = useCallback(
    (tokens: AuthTokens | null) => {
      clearRefreshTimer();
      if (!tokens) {
        return;
      }
      const msUntilRefresh = Math.max(tokens.expiresAt - Date.now() - 60_000, 5_000);
      refreshTimer.current = window.setTimeout(async () => {
        try {
          const next = await refreshTokensRequest(tokens.refreshToken);
          const mapped = mapTokens(next, state.tenantId, state.userId);
          setState(mapped);
          persistState(mapped);
        } catch (error) {
          console.error('Erro ao renovar tokens, realizando logout.', error);
          setState(DEFAULT_STATE);
          persistState(DEFAULT_STATE);
        }
      }, msUntilRefresh);
    },
    [state.tenantId, state.userId],
  );

  useEffect(() => {
    setInitializing(false);
  }, []);

  useEffect(() => {
    scheduleRefresh(state.tokens);
    return clearRefreshTimer;
  }, [state.tokens, scheduleRefresh]);

  const login = useCallback(
    async (credentials: LoginRequest) => {
      const payload = await loginRequest(credentials);
      const mapped = mapTokens(payload);
      setState(mapped);
      persistState(mapped);
    },
    [],
  );

  const signup = useCallback(async (input: SignupRequest) => {
    const result = await signupRequest(input);
    const mapped = mapSignupResult(result);
    setState(mapped);
    persistState(mapped);
  }, []);

  const logout = useCallback(() => {
    clearRefreshTimer();
    setState(DEFAULT_STATE);
    persistState(DEFAULT_STATE);
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      isAuthenticated: Boolean(state.tokens),
      initializing,
      tenantId: state.tenantId,
      userId: state.userId,
      accessToken: state.tokens?.accessToken ?? null,
      login,
      signup,
      logout,
    }),
    [state, initializing, login, signup, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = (): AuthContextValue => {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth deve ser usado dentro de AuthProvider');
  }
  return ctx;
};
