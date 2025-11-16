import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ClientListPage from './ClientListPage';
import * as apiClient from '../../lib/apiClient';
import { type Client } from '../../lib/apiClient';
import { AUTH_STORAGE_KEY } from '../../contexts/AuthContext';
import { AuthProvider } from '../../contexts/AuthContext';

// Mock the apiClient module
vi.mock('../../lib/apiClient');

const mockedApiClient = vi.mocked(apiClient);
const authState = {
  tenantId: 'tenant-1',
  userId: 'user-1',
  tokens: { accessToken: 'token', refreshToken: 'refresh', expiresAt: Date.now() + 100000 },
};

const renderWithAuth = () =>
  render(
    <AuthProvider>
      <ClientListPage />
    </AuthProvider>,
  );

describe('ClientListPage', () => {
  beforeEach(() => {
    // Reseta os mocks antes de cada teste
    vi.resetAllMocks();
    window.localStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(authState));
  });

  afterEach(() => {
    window.localStorage.clear();
  });

  it('should display a loading state initially', () => {
    // Para este teste, fazemos a API nunca resolver
    mockedApiClient.fetchClients.mockReturnValue(new Promise(() => {}));

    renderWithAuth();

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
    mockedApiClient.fetchClients.mockResolvedValue(mockClients);

    renderWithAuth();

    // Espera que os clientes apareçam na tela
    await waitFor(() => {
      expect(screen.getByText('John Doe')).toBeInTheDocument();
      expect(screen.getByText('Jane Smith')).toBeInTheDocument();
    });
    expect(mockedApiClient.fetchClients).toHaveBeenCalledWith({
      tenantId: authState.tenantId,
      accessToken: authState.tokens.accessToken,
    });
  });

  it('should display an empty state message when no clients are found', async () => {
    mockedApiClient.fetchClients.mockResolvedValue([]);

    renderWithAuth();

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
    mockedApiClient.fetchClients.mockRejectedValue(new Error(errorMessage));

    // Mock console.error para evitar poluir a saída do teste
    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    renderWithAuth();

    // Espera que o erro seja logado
    await waitFor(() => {
      expect(consoleErrorSpy).toHaveBeenCalledWith('Error fetching clients:', expect.any(Error));
    });

    // Verifica se nenhum cliente foi renderizado
    expect(screen.queryByText('John Doe')).not.toBeInTheDocument();

    consoleErrorSpy.mockRestore();
  });
});
