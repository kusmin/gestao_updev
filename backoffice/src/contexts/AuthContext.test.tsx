import { act, renderHook, waitFor } from '@testing-library/react';
import { type FC, type ReactNode } from 'react';
import { vi } from 'vitest';
import { AuthProvider, useAuth } from './AuthContext';
import * as apiClient from '../lib/apiClient';

describe('AuthContext', () => {
  let setAuthTokenSpy: ReturnType<typeof vi.spyOn>;

  const wrapper: FC<{ children: ReactNode }> = ({ children }) => (
    <AuthProvider>{children}</AuthProvider>
  );

  beforeEach(() => {
    localStorage.clear();
    setAuthTokenSpy = vi.spyOn(apiClient, 'setAuthToken');
  });

  afterEach(() => {
    setAuthTokenSpy.mockRestore();
  });

  it('stores token on login and updates api client', () => {
    const { result } = renderHook(() => useAuth(), { wrapper });

    act(() => {
      result.current.login('test-token');
    });

    expect(result.current.token).toBe('test-token');
    expect(localStorage.getItem('backoffice_token')).toBe('test-token');
    expect(setAuthTokenSpy).toHaveBeenCalledWith('test-token');
  });

  it('clears token on logout and updates api client', () => {
    const { result } = renderHook(() => useAuth(), { wrapper });

    act(() => {
      result.current.login('test-token');
      result.current.logout();
    });

    expect(result.current.token).toBeNull();
    expect(localStorage.getItem('backoffice_token')).toBeNull();
    expect(setAuthTokenSpy).toHaveBeenLastCalledWith(null);
  });

  it('initializes token from localStorage on mount', async () => {
    localStorage.setItem('backoffice_token', 'existing-token');

    const { result } = renderHook(() => useAuth(), { wrapper });

    await waitFor(() => {
      expect(result.current.token).toBe('existing-token');
    });
    expect(setAuthTokenSpy).toHaveBeenCalledWith('existing-token');
  });
});
