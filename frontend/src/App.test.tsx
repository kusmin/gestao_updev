import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import App from './App';
import { beforeAll, beforeEach, afterEach, describe, it, expect, vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AuthState } from './contexts/AuthContextDefinition';

// Mock the apiClient used in ClientListPage
vi.mock('./lib/apiClient', () => ({
  fetchClients: vi.fn(() => Promise.resolve([])),
}));

const MOCK_AUTH_STATE: AuthState = {
  tenantId: 'tenant-1',
  userId: 'user-1',
  tokens: {
    accessToken: 'fake-token',
    refreshToken: 'refresh-token',
    expiresAt: Date.now() + 60 * 60 * 1000,
  },
};

// Mock authUtils to control the initial state of AuthProvider
vi.mock('./contexts/authUtils', () => ({
  loadStoredAuth: vi.fn(() => MOCK_AUTH_STATE),
  persistState: vi.fn(),
  DEFAULT_STATE: {
    tenantId: null,
    userId: null,
    tokens: null,
  },
  mapSignupResult: vi.fn(),
  mapTokens: vi.fn(),
}));

const queryClient = new QueryClient();

beforeAll(() => {
  vi.stubGlobal('matchMedia', (query: string) => ({
    matches: query.includes('prefers-color-scheme') ? false : true,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  }));
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe('App', () => {
  it('should navigate to the clients page from the navigation menu', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>,
    );

    const navLink = await screen.findByRole('link', { name: /clientes/i });
    await userEvent.click(navLink);

    const heading = await screen.findByRole('heading', { name: /clientes/i });
    expect(heading).toBeInTheDocument();
  });
});
