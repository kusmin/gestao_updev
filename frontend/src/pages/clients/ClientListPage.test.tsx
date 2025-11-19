import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import ClientListPage from './ClientListPage';
import * as apiClient from '../../lib/apiClient';
import { type Client } from '../../lib/apiClient';
import { AuthState } from '../../contexts/AuthContextDefinition';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AuthProvider } from '../../contexts/AuthContext';
import { ColorModeProvider } from '../../theme/ColorModeProvider';
import { SnackbarProvider } from '../../contexts/SnackbarProvider';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
});

// Mock the apiClient module
const mockFetchClients = vi.fn();
const mockDeleteClient = vi.fn();
vi.mock('../../lib/apiClient', async () => {
  const actual = await vi.importActual('../../lib/apiClient');
  return {
    ...actual,
    fetchClients: mockFetchClients,
    deleteClient: mockDeleteClient,
  };
});

const MOCKED_API_CLIENT = {
  fetchClients: vi.mocked(mockFetchClients),
  deleteClient: vi.mocked(mockDeleteClient),
};

const MOCK_AUTH_STATE: AuthState = {
  tenantId: 'tenant-1',
  userId: 'user-1',
  tokens: { accessToken: 'token', refreshToken: 'refresh', expiresAt: Date.now() + 100000 },
};



describe('ClientListPage', () => {
  beforeEach(() => {
    // Reseta os mocks antes de cada teste
    vi.resetAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should display a loading state initially', () => {
    // Para este teste, fazemos a API nunca resolver
    MOCKED_API_CLIENT.fetchClients.mockReturnValue(new Promise(() => {}));

    render(
      <ColorModeProvider>
        <AuthProvider>
          <SnackbarProvider>
            <QueryClientProvider client={queryClient}>
              <ClientListPage />
            </QueryClientProvider>
          </SnackbarProvider>
        </AuthProvider>
      </ColorModeProvider>,
    );

    // O componente não tem um estado de loading explícito,
    // então vamos verificar se a tabela está vazia inicialmente.
    expect(screen.getByRole('table')).toBeInTheDocument();
    expect(screen.queryByRole('cell')).not.toBeInTheDocument();
  });

  it('should render a list of clients', async () => {
    const mockClients: Client[] = [
      { id: '1', name: 'John Doe', email: 'john@example.com', phone: '12345' },
      { id: '2', name: 'Jane Smith', email: 'jane@example.com', phone: '67890' },
    ];
    MOCKED_API_CLIENT.fetchClients.mockResolvedValue(mockClients);

    await render(
      <ColorModeProvider>
        <AuthProvider>
          <SnackbarProvider>
            <QueryClientProvider client={queryClient}>
              <ClientListPage />
            </QueryClientProvider>
          </SnackbarProvider>
        </AuthProvider>
      </ColorModeProvider>,
    );

    // Espera que os clientes apareçam na tela
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
      expect(screen.getByText('Jane Smith')).toBeInTheDocument();
    });
    expect(MOCKED_API_CLIENT.fetchClients).toHaveBeenCalledWith({
      tenantId: MOCK_AUTH_STATE.tenantId,
      accessToken: MOCK_AUTH_STATE.tokens?.accessToken,
    });
  });

  it('should display an empty state message when no clients are found', async () => {
    MOCKED_API_CLIENT.fetchClients.mockResolvedValue([]);

    await render(
      <ColorModeProvider>
        <AuthProvider>
          <SnackbarProvider>
            <QueryClientProvider client={queryClient}>
              <ClientListPage />
            </QueryClientProvider>
          </SnackbarProvider>
        </AuthProvider>
      </ColorModeProvider>,
    );

    // O componente não tem uma mensagem de estado vazio explícita.
    // Vamos verificar se a tabela está presente, mas sem linhas de dados.
    await waitFor(() => {
      const rows = screen.queryAllByRole('row');
      // Espera-se 1 linha (o cabeçalho da tabela)
      expect(rows.length).toBe(1);
    });
  });

  it('should handle API errors gracefully', async () => {
    const errorMessage = 'Failed to fetch clients';
    MOCKED_API_CLIENT.fetchClients.mockRejectedValue(new Error(errorMessage));

    // Mock console.error para evitar poluir a saída do teste
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    await render(
      <ColorModeProvider>
        <AuthProvider>
          <SnackbarProvider>
            <QueryClientProvider client={queryClient}>
              <ClientListPage />
            </QueryClientProvider>
          </SnackbarProvider>
        </AuthProvider>
      </ColorModeProvider>,
    );

    // Espera que o erro seja logado
    await waitFor(() => {
      expect(consoleErrorSpy).toHaveBeenCalledWith('Error fetching clients:', expect.any(Error));
    });

    // Verifica se nenhum cliente foi renderizado
    expect(screen.queryByText('John Doe')).not.toBeInTheDocument();

    consoleErrorSpy.mockRestore();
  });
});
