import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import ClientForm from './ClientForm';
import { createClient, updateClient } from '../../lib/apiClient';

vi.mock('../../lib/apiClient', () => ({
  createClient: vi.fn(),
  updateClient: vi.fn(),
}));

const mockCreate = vi.mocked(createClient);
const mockUpdate = vi.mocked(updateClient);
const TENANT_ID = 'a4b2b2b2-b2b2-4b2b-b2b2-b2b2b2b2b2b2';

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

  it('cria um novo cliente quando client é nulo', async () => {
    mockCreate.mockResolvedValueOnce({
      id: '1',
      name: 'Alice',
      email: 'alice@example.com',
      phone: '123',
    });

    render(<ClientForm {...baseProps} client={null} />);

    await userEvent.type(screen.getByLabelText(/Nome/i), 'Alice');
    await userEvent.type(screen.getByLabelText(/Email/i), 'alice@example.com');
    await userEvent.type(screen.getByLabelText(/Telefone/i), '123');
    await userEvent.click(screen.getByRole('button', { name: /Salvar/i }));

    await waitFor(() =>
      expect(mockCreate).toHaveBeenCalledWith(TENANT_ID, {
        name: 'Alice',
        email: 'alice@example.com',
        phone: '123',
      }),
    );
    expect(baseProps.onSave).toHaveBeenCalledWith(
      expect.objectContaining({ id: '1', name: 'Alice' }),
    );
    expect(baseProps.onClose).toHaveBeenCalled();
  });

  it('atualiza um cliente existente quando client é fornecido', async () => {
    const existing = { id: '42', name: 'Bob', email: 'bob@old.com', phone: '000' };
    mockUpdate.mockResolvedValueOnce({
      ...existing,
      name: 'Bob Atualizado',
      email: 'bob@new.com',
    });

    render(<ClientForm {...baseProps} client={existing} />);

    expect(screen.getByDisplayValue('Bob')).toBeInTheDocument();
    await userEvent.clear(screen.getByLabelText(/Nome/i));
    await userEvent.type(screen.getByLabelText(/Nome/i), 'Bob Atualizado');
    await userEvent.clear(screen.getByLabelText(/Email/i));
    await userEvent.type(screen.getByLabelText(/Email/i), 'bob@new.com');
    await userEvent.click(screen.getByRole('button', { name: /Salvar/i }));

    await waitFor(() =>
      expect(mockUpdate).toHaveBeenCalledWith(TENANT_ID, '42', {
        name: 'Bob Atualizado',
        email: 'bob@new.com',
        phone: '000',
      }),
    );
    expect(baseProps.onSave).toHaveBeenCalledWith(
      expect.objectContaining({ id: '42', name: 'Bob Atualizado' }),
    );
    expect(baseProps.onClose).toHaveBeenCalled();
  });
});
