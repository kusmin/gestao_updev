import type { components } from '../types/api';

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/v1';

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

// fetchHealth consulta o endpoint público /v1/healthz e retorna os dados relevantes.
export const fetchHealth = async (): Promise<HealthResponse> => {
  const res = await fetch(`${API_BASE.replace(/\/$/, '')}/healthz`);
  if (!res.ok) {
    throw new Error(`Falha ao consultar healthz: ${res.status}`);
  }
  const payload = (await res.json()) as {
    data: HealthResponse;
  };
  return payload.data;
};

export const fetchClients = async (tenantId: string): Promise<Client[]> => {
  const res = await fetch(`${API_BASE.replace(/\/$/, '')}/clients`, {
    headers: {
      'X-Tenant-ID': tenantId,
    },
  });
  if (!res.ok) {
    throw new Error(`Falha ao consultar clientes: ${res.status}`);
  }
  const payload = (await res.json()) as {
    data?: Client[];
  };
  return payload.data ?? [];
};

export const deleteClient = async (tenantId: string, clientId: string): Promise<void> => {
  const res = await fetch(`${API_BASE.replace(/\/$/, '')}/clients/${clientId}`, {
    method: 'DELETE',
    headers: {
      'X-Tenant-ID': tenantId,
    },
  });
  if (!res.ok) {
    throw new Error(`Falha ao remover cliente: ${res.status}`);
  }
};

export const createClient = async (tenantId: string, client: ClientRequest): Promise<Client> => {
  const res = await fetch(`${API_BASE.replace(/\/$/, '')}/clients`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Tenant-ID': tenantId,
    },
    body: JSON.stringify(client),
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
  tenantId: string,
  clientId: string,
  client: ClientRequest,
): Promise<Client> => {
  const res = await fetch(`${API_BASE.replace(/\/$/, '')}/clients/${clientId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'X-Tenant-ID': tenantId,
    },
    body: JSON.stringify(client),
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
