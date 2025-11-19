import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import ClientForm from './ClientForm';
import { createClient, updateClient } from '../../lib/apiClient';
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
const mockCreate = vi.fn();
const mockUpdate = vi.fn();
vi.mock('../../lib/apiClient', async () => {
  const actual = await vi.importActual('../../lib/apiClient');
  return {
    ...actual,
    createClient: mockCreate,
    updateClient: mockUpdate,
  };
});

const TENANT_ID = 'tenant-ctx';
const ACCESS_TOKEN = 'token-ctx';
const MOCK_AUTH_STATE: AuthState = {
  tenantId: TENANT_ID,
  userId: 'user-ctx',
  tokens: { accessToken: ACCESS_TOKEN, refreshToken: 'refresh', expiresAt: Date.now() + 100000 },
};

describe('ClientForm', () => {
  const baseProps = {
    open: true,
    onClose: vi.fn(),
    onSave: vi.fn(),
  };

  beforeEach(() => {
    vi.clearAllMocks();
    baseProps.onClose.mockReset();
    baseProps.onSave.mockReset();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('cria um novo cliente quando client é nulo', async () => {
    createClient.mockResolvedValueOnce({
      id: '1',
      name: 'Alice',
      email: 'alice@example.com',
      phone: '123',
    });

    render(
      <ColorModeProvider>
        <AuthProvider>
          <SnackbarProvider>
            <QueryClientProvider client={queryClient}>
              <ClientForm {...baseProps} client={null} />
            </QueryClientProvider>
          </SnackbarProvider>
        </AuthProvider>
      </ColorModeProvider>,
    );

    await userEvent.type(screen.getByLabelText(/Nome/i), 'Alice');
    await userEvent.type(screen.getByLabelText(/Email/i), 'alice@example.com');
    await userEvent.type(screen.getByLabelText(/Telefone/i), '123');
    await userEvent.click(screen.getByRole('button', { name: /Salvar/i }));

    await waitFor(() =>
      expect(createClient).toHaveBeenCalledWith({
        tenantId: TENANT_ID,
        accessToken: ACCESS_TOKEN,
        input: {
          name: 'Alice',
          email: 'alice@example.com',
          phone: '123',
        },
      }),
    );
    expect(baseProps.onSave).toHaveBeenCalledWith(
      expect.objectContaining({ id: '1', name: 'Alice' }),
    );
    expect(baseProps.onClose).toHaveBeenCalled();
  });

  it('atualiza um cliente existente quando client é fornecido', async () => {
    const existing = { id: '42', name: 'Bob', email: 'bob@old.com', phone: '000' };
    updateClient.mockResolvedValueOnce({
      ...existing,
      name: 'Bob Atualizado',
      email: 'bob@new.com',
    });

    render(
      <ColorModeProvider>
        <AuthProvider>
          <SnackbarProvider>
            <QueryClientProvider client={queryClient}>
              <ClientForm {...baseProps} client={existing} />
            </QueryClientProvider>
          </SnackbarProvider>
        </AuthProvider>
      </ColorModeProvider>,
    );

    expect(screen.getByDisplayValue('Bob')).toBeInTheDocument();
    await userEvent.clear(screen.getByLabelText(/Nome/i));
    await userEvent.type(screen.getByLabelText(/Nome/i), 'Bob Atualizado');
    await userEvent.clear(screen.getByLabelText(/Email/i));
    await userEvent.type(screen.getByLabelText(/Email/i), 'bob@new.com');
    await userEvent.click(screen.getByRole('button', { name: /Salvar/i }));

    await waitFor(() =>
      expect(updateClient).toHaveBeenCalledWith({
        tenantId: TENANT_ID,
        clientId: '42',
        accessToken: ACCESS_TOKEN,
        input: {
          name: 'Bob Atualizado',
          email: 'bob@new.com',
          phone: '000',
        },
      }),
    );
    expect(baseProps.onSave).toHaveBeenCalledWith(
      expect.objectContaining({ id: '42', name: 'Bob Atualizado' }),
    );
    expect(baseProps.onClose).toHaveBeenCalled();
  });
});
