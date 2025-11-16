import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import App from './App';
import { beforeAll, describe, it, expect, vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

// Mock the apiClient used in ClientListPage
vi.mock('./lib/apiClient', () => ({
  fetchClients: vi.fn(() => Promise.resolve([])),
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
