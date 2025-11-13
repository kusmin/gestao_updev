import { render, screen } from '@testing-library/react';
import App from './App';
import { describe, it, expect } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

// Mock the apiClient used in ClientListPage
vi.mock('./lib/apiClient', () => ({
  fetchClients: vi.fn(() => Promise.resolve([])),
}));

const queryClient = new QueryClient();

describe('App', () => {
  it('should render the clients list page', () => {
    render(
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    );

    expect(screen.getByRole('heading', { name: /clientes/i })).toBeInTheDocument();
  });
});
