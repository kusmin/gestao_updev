import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import App from './App';
import { beforeAll, beforeEach, afterEach, describe, it, expect, vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AUTH_STORAGE_KEY } from './contexts/AuthContext';

// Mock the apiClient used in ClientListPage
vi.mock('./lib/apiClient', () => ({
  fetchClients: vi.fn(() => Promise.resolve([])),
}));

const queryClient = new QueryClient();
const MOCK_AUTH_STATE = {
  tenantId: 'tenant-1',
  userId: 'user-1',
  tokens: {
    accessToken: 'fake-token',
    refreshToken: 'refresh-token',
    expiresAt: Date.now() + 60 * 60 * 1000,
  },
};

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

beforeEach(() => {
  window.localStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(MOCK_AUTH_STATE));
});

afterEach(() => {
  window.localStorage.clear();
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
