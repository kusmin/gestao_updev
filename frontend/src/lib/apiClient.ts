import type { components } from '../types/api';

const API_BASE = (import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/v1').replace(/\/$/, '');

type HealthResponse = {
  status?: string;
  env?: string;
};

export interface Client {
  id: string;
  name: string;
  email?: string | null;
  phone?: string | null;
  notes?: string | null;
}

export type ClientRequest = components['schemas']['handler.ClientRequest'];
export type LoginRequest = components['schemas']['handler.LoginRequest'];
export type SignupRequest = components['schemas']['handler.SignupRequest'];
export type SignupResponse = components['schemas']['handler.SignupResponse'];

export type AuthTokensResponse = {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  tenant_id?: string;
  user_id?: string;
};

type AuthenticatedRequest = {
  tenantId: string;
  accessToken: string;
};

const withTenantHeaders = (params: AuthenticatedRequest) => ({
  'X-Tenant-ID': params.tenantId,
  Authorization: `Bearer ${params.accessToken}`,
});

const jsonHeaders = {
  'Content-Type': 'application/json',
};

// fetchHealth consulta o endpoint público /v1/healthz e retorna os dados relevantes.
export const fetchHealth = async (): Promise<HealthResponse> => {
  const res = await fetch(`${API_BASE}/healthz`);
  if (!res.ok) {
    throw new Error(`Falha ao consultar healthz: ${res.status}`);
  }
  const payload = (await res.json()) as {
    data: HealthResponse;
  };
  return payload.data;
};

export const fetchClients = async (params: AuthenticatedRequest): Promise<Client[]> => {
  const res = await fetch(`${API_BASE}/clients`, {
    headers: withTenantHeaders(params),
  });
  if (!res.ok) {
    throw new Error(`Falha ao consultar clientes: ${res.status}`);
  }
  const payload = (await res.json()) as {
    data?: Client[];
  };
  return payload.data ?? [];
};

export const deleteClient = async (
  params: AuthenticatedRequest & { clientId: string },
): Promise<void> => {
  const res = await fetch(`${API_BASE}/clients/${params.clientId}`, {
    method: 'DELETE',
    headers: withTenantHeaders(params),
  });
  if (!res.ok) {
    throw new Error(`Falha ao remover cliente: ${res.status}`);
  }
};

export const createClient = async (
  params: AuthenticatedRequest & { input: ClientRequest },
): Promise<Client> => {
  const res = await fetch(`${API_BASE}/clients`, {
    method: 'POST',
    headers: {
      ...jsonHeaders,
      ...withTenantHeaders(params),
    },
    body: JSON.stringify(params.input),
  });
  if (!res.ok) {
    throw new Error(`Falha ao criar cliente: ${res.status}`);
  }
  const payload = (await res.json()) as { data?: Client };
  if (!payload.data) {
    throw new Error('Resposta inesperada ao criar cliente');
  }
  return payload.data;
};

export const updateClient = async (
  params: AuthenticatedRequest & { clientId: string; input: ClientRequest },
): Promise<Client> => {
  const res = await fetch(`${API_BASE}/clients/${params.clientId}`, {
    method: 'PUT',
    headers: {
      ...jsonHeaders,
      ...withTenantHeaders(params),
    },
    body: JSON.stringify(params.input),
  });
  if (!res.ok) {
    throw new Error(`Falha ao atualizar cliente: ${res.status}`);
  }
  const payload = (await res.json()) as { data?: Client };
  if (!payload.data) {
    throw new Error('Resposta inesperada ao atualizar cliente');
  }
  return payload.data;
};

export const login = async (credentials: LoginRequest): Promise<AuthTokensResponse> => {
  const res = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: jsonHeaders,
    body: JSON.stringify(credentials),
  });
  if (!res.ok) {
    throw new Error(`Falha ao autenticar: ${res.status}`);
  }
  const payload = (await res.json()) as { data?: AuthTokensResponse };
  if (!payload.data) {
    throw new Error('Resposta inesperada ao autenticar');
  }
  return payload.data;
};

export const signup = async (input: SignupRequest): Promise<SignupResponse> => {
  const res = await fetch(`${API_BASE}/auth/signup`, {
    method: 'POST',
    headers: jsonHeaders,
    body: JSON.stringify(input),
  });
  if (!res.ok) {
    throw new Error(`Falha ao criar conta: ${res.status}`);
  }
  const payload = (await res.json()) as { data?: SignupResponse };
  if (!payload.data) {
    throw new Error('Resposta inesperada ao criar conta');
  }
  return payload.data;
};

export const refreshTokens = async (refreshToken: string): Promise<AuthTokensResponse> => {
  const res = await fetch(`${API_BASE}/auth/refresh`, {
    method: 'POST',
    headers: jsonHeaders,
    body: JSON.stringify({ refresh_token: refreshToken }),
  });
  if (!res.ok) {
    throw new Error(`Falha ao renovar tokens: ${res.status}`);
  }
  const payload = (await res.json()) as { data?: AuthTokensResponse };
  if (!payload.data) {
    throw new Error('Resposta inesperada ao renovar tokens');
  }
  return payload.data;
};
