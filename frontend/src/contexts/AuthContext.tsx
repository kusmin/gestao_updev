import { AuthContext, AuthContextValue, AuthState, AuthTokens } from './AuthContextDefinition';
export { AuthContext } from './AuthContextDefinition';
export type { AuthContextValue } from './AuthContextDefinition';
import {
  DEFAULT_STATE,
  loadStoredAuth,
  mapSignupResult,
  mapTokens,
  persistState,
} from './authUtils';
import {
  login as loginRequest,
  signup as signupRequest,
  refreshTokens as refreshTokensRequest,
  LoginRequest,
  SignupRequest,
} from '../lib/apiClient';
import {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react';

type AuthProviderProps = {
  children: React.ReactNode;
};

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [state, setState] = useState<AuthState>(() => loadStoredAuth());
  const [initializing, setInitializing] = useState(true);
  const refreshTimer = useRef<number | null>(null);

  const clearRefreshTimer = useCallback(() => {
    if (refreshTimer.current) {
      window.clearTimeout(refreshTimer.current);
      refreshTimer.current = null;
    }
  }, []);

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
    [state.tenantId, state.userId, clearRefreshTimer],
  );

  useEffect(() => {
    setInitializing(false);
  }, []);

  useEffect(() => {
    scheduleRefresh(state.tokens);
    return clearRefreshTimer;
  }, [state.tokens, scheduleRefresh, clearRefreshTimer]);

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
  }, [clearRefreshTimer]);

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


