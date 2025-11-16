import { describe, it, expect, vi, beforeEach, afterEach, afterAll } from 'vitest';
import {
  fetchHealth,
  fetchClients,
  createClient,
  updateClient,
  deleteClient,
  type ClientRequest,
  login,
  signup,
  refreshTokens,
  type SignupRequest,
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

    const result = await fetchClients({ tenantId: 'tenant-1', accessToken: 'token' });
    expect(result).toEqual([]);
    expect(fetchMock()).toHaveBeenCalledWith('http://localhost:8080/v1/clients', {
      headers: { 'X-Tenant-ID': 'tenant-1', Authorization: 'Bearer token' },
    });
  });

  it('createClient envia POST e retorna cliente criado', async () => {
    const payload: ClientRequest = { name: 'Alice', email: 'alice@example.com', phone: '123' };
    fetchMock().mockResolvedValueOnce({
      ok: true,
      status: 201,
      json: async () => ({ data: { id: '1', ...payload } }),
    } as Response);

    const result = await createClient({ tenantId: 'tenant-2', accessToken: 'token', input: payload });
    expect(result).toMatchObject({ id: '1', name: 'Alice' });
    expect(fetchMock()).toHaveBeenCalledWith('http://localhost:8080/v1/clients', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Tenant-ID': 'tenant-2',
        Authorization: 'Bearer token',
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

    const result = await updateClient({
      tenantId: 'tenant-3',
      clientId: '42',
      accessToken: 'token',
      input: payload,
    });
    expect(result).toMatchObject({ id: '42', name: 'Bob' });
    expect(fetchMock()).toHaveBeenCalledWith('http://localhost:8080/v1/clients/42', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'X-Tenant-ID': 'tenant-3',
        Authorization: 'Bearer token',
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

    await deleteClient({ tenantId: 'tenant-4', clientId: '77', accessToken: 'token' });
    expect(fetchMock()).toHaveBeenCalledWith('http://localhost:8080/v1/clients/77', {
      method: 'DELETE',
      headers: { 'X-Tenant-ID': 'tenant-4', Authorization: 'Bearer token' },
    });
  });

  it('login retorna tokens quando sucesso', async () => {
    fetchMock().mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({
        data: { access_token: 'at', refresh_token: 'rt', expires_in: 900 },
      }),
    } as Response);

    const result = await login({ email: 'foo@bar.com', password: 'secret' });
    expect(result).toMatchObject({ access_token: 'at', refresh_token: 'rt' });
    expect(fetchMock()).toHaveBeenCalledWith('http://localhost:8080/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: 'foo@bar.com', password: 'secret' }),
    });
  });

  it('signup retorna dados do tenant', async () => {
    const payload: SignupRequest = {
      company: { name: 'Empresa', document: '123', phone: '111' },
      user: { name: 'Admin', email: 'admin@test.com', password: 'Senha@123', phone: '222' },
    };
    fetchMock().mockResolvedValueOnce({
      ok: true,
      status: 201,
      json: async () => ({
        data: {
          tenant_id: 'tenant-1',
          user_id: 'user-1',
          access_token: 'at',
          refresh_token: 'rt',
          expires_in: 900,
        },
      }),
    } as Response);

    const result = await signup(payload);
    expect(result).toMatchObject({ tenant_id: 'tenant-1', user_id: 'user-1' });
  });

  it('refreshTokens devolve novo par de tokens', async () => {
    fetchMock().mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => ({
        data: { access_token: 'new-at', refresh_token: 'new-rt', expires_in: 900 },
      }),
    } as Response);

    const result = await refreshTokens('old-rt');
    expect(result.access_token).toBe('new-at');
    expect(fetchMock()).toHaveBeenCalledWith('http://localhost:8080/v1/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: 'old-rt' }),
    });
  });
});
