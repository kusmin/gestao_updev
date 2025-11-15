import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import ResourceListPage from './index';
import apiClient from '@/lib/apiClient';

vi.mock('@/lib/apiClient', () => ({
  default: vi.fn(),
}));

type ClientResource = {
  id: string;
  name: string;
  email: string;
};

const mockApi = vi.mocked(apiClient);

const TestForm = ({
  open,
  onClose,
  onSave,
  item,
}: {
  open: boolean;
  onClose: () => void;
  onSave: (item: Partial<ClientResource>) => Promise<void>;
  item: ClientResource | null;
}) => {
  if (!open) return null;

  const handleSave = async () => {
    const payload =
      item === null
        ? { name: 'New Client', email: 'new@example.com' }
        : { name: `${item.name} Updated`, email: item.email };

    await onSave(payload);
    onClose();
  };

  return (
    <div data-testid="resource-form">
      <span data-testid="form-mode">{item ? `editing-${item.name}` : 'creating'}</span>
      <button type="button" onClick={handleSave}>
        Salvar
      </button>
    </div>
  );
};

describe('ResourceListPage', () => {
  beforeEach(() => {
    mockApi.mockReset();
  });

  it('lists, creates, edits and deletes resources', async () => {
    const initialItems: ClientResource[] = [{ id: '1', name: 'Alice', email: 'alice@example.com' }];
    const afterCreate = [
      ...initialItems,
      { id: '2', name: 'New Client', email: 'new@example.com' },
    ];
    const afterUpdate = [
      { id: '1', name: 'Alice Updated', email: 'alice@example.com' },
      afterCreate[1],
    ];
    const afterDelete = [afterCreate[1]];

    mockApi
      .mockResolvedValueOnce({ data: initialItems }) // initial fetch
      .mockResolvedValueOnce({}) // create
      .mockResolvedValueOnce({ data: afterCreate }) // refetch after create
      .mockResolvedValueOnce({}) // update
      .mockResolvedValueOnce({ data: afterUpdate }) // refetch after update
      .mockResolvedValueOnce({}) // delete
      .mockResolvedValueOnce({ data: afterDelete }); // refetch after delete

    render(
      <ResourceListPage<ClientResource>
        title="Clientes"
        endpoint="/clients"
        columns={[
          { header: 'Nome', accessor: 'name' },
          { header: 'Email', accessor: 'email' },
        ]}
        renderForm={(props) => <TestForm {...props} />}
      />,
    );

    // initial fetch
    expect(await screen.findByText('Alice')).toBeInTheDocument();

    // create flow
    await userEvent.click(screen.getByRole('button', { name: /Adicionar Clientes/i }));
    expect(screen.getByTestId('form-mode').textContent).toBe('creating');
    await userEvent.click(screen.getByRole('button', { name: /Salvar/i }));

    await waitFor(() =>
      expect(mockApi).toHaveBeenCalledWith(
        '/clients',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ name: 'New Client', email: 'new@example.com' }),
        }),
      ),
    );
    expect(await screen.findByText('New Client')).toBeInTheDocument();

    // edit flow
    await userEvent.click(screen.getAllByText('Editar')[1]);
    expect(screen.getByTestId('form-mode').textContent).toBe('editing-New Client');
    await userEvent.click(screen.getByRole('button', { name: /Salvar/i }));

    await waitFor(() =>
      expect(mockApi).toHaveBeenCalledWith(
        '/clients/2',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify({ name: 'New Client Updated', email: 'new@example.com' }),
        }),
      ),
    );
    expect(await screen.findByText('Alice Updated')).toBeInTheDocument();

    // delete flow
    await userEvent.click(screen.getAllByText('Excluir')[0]);
    await waitFor(() =>
      expect(mockApi).toHaveBeenCalledWith(
        '/clients/1',
        expect.objectContaining({
          method: 'DELETE',
        }),
      ),
    );
    await waitFor(() => expect(screen.queryByText('Alice Updated')).not.toBeInTheDocument());
    expect(screen.getByText('New Client')).toBeInTheDocument();
  });
});
