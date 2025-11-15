import { describe, it, expect, vi, beforeEach, afterEach, afterAll } from 'vitest';
import {
  fetchHealth,
  fetchClients,
  createClient,
  updateClient,
  deleteClient,
  type ClientRequest,
} from './apiClient';

const originalFetch = global.fetch;

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn());
});

afterEach(() => {
  vi.unstubAllGlobals();
});

afterAll(() => {
  global.fetch = originalFetch;
});

const fetchMock = () => global.fetch as unknown as ReturnType<typeof vi.fn>;

describe('apiClient', () => {
  it('fetchHealth retorna payload', async () => {
    fetchMock().mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({ data: { status: 'ok', env: 'test' } }),
    } as Response);

    const data = await fetchHealth();
    expect(fetchMock()).toHaveBeenCalledWith('http://localhost:8080/v1/healthz');
    expect(data).toEqual({ status: 'ok', env: 'test' });
  });

  it('fetchHealth lança erro quando resposta não é ok', async () => {
    fetchMock().mockResolvedValueOnce({
      ok: false,
      status: 503,
      json: async () => ({}),
    } as Response);

    await expect(fetchHealth()).rejects.toThrow(/503/);
  });

  it('fetchClients retorna lista vazia quando payload não contém data', async () => {
    fetchMock().mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({}),
    } as Response);

    const result = await fetchClients('tenant-1');
    expect(result).toEqual([]);
    expect(fetchMock()).toHaveBeenCalledWith('http://localhost:8080/v1/clients', {
      headers: { 'X-Tenant-ID': 'tenant-1' },
    });
  });

  it('createClient envia POST e retorna cliente criado', async () => {
    const payload: ClientRequest = { name: 'Alice', email: 'alice@example.com', phone: '123' };
    fetchMock().mockResolvedValueOnce({
      ok: true,
      status: 201,
      json: async () => ({ data: { id: '1', ...payload } }),
    } as Response);

    const result = await createClient('tenant-2', payload);
    expect(result).toMatchObject({ id: '1', name: 'Alice' });
    expect(fetchMock()).toHaveBeenCalledWith('http://localhost:8080/v1/clients', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Tenant-ID': 'tenant-2',
      },
      body: JSON.stringify(payload),
    });
  });

  it('updateClient envia PUT e retorna cliente atualizado', async () => {
    const payload: ClientRequest = { name: 'Bob', email: 'bob@example.com', phone: '' };
    fetchMock().mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({ data: { id: '42', ...payload } }),
    } as Response);

    const result = await updateClient('tenant-3', '42', payload);
    expect(result).toMatchObject({ id: '42', name: 'Bob' });
    expect(fetchMock()).toHaveBeenCalledWith('http://localhost:8080/v1/clients/42', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'X-Tenant-ID': 'tenant-3',
      },
      body: JSON.stringify(payload),
    });
  });

  it('deleteClient envia DELETE', async () => {
    fetchMock().mockResolvedValueOnce({
      ok: true,
      status: 204,
      json: async () => ({}),
    } as Response);

    await deleteClient('tenant-4', '77');
    expect(fetchMock()).toHaveBeenCalledWith('http://localhost:8080/v1/clients/77', {
      method: 'DELETE',
      headers: { 'X-Tenant-ID': 'tenant-4' },
    });
  });
});
