import { render, screen } from '@testing-library/react';
import App from './App';
import { describe, it, expect, vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

// Mock the apiClient used in ClientListPage
vi.mock('./lib/apiClient', () => ({
  fetchClients: vi.fn(() => Promise.resolve([])),
}));

const queryClient = new QueryClient();

describe('App', () => {
  it('should render the clients list page', async () => {
    render(
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    );

    const heading = await screen.findByRole('heading', { name: /clientes/i });
    expect(heading).toBeInTheDocument();
  });
});
