import type { paths, components } from '../types/api';

const API_BASE =
  import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/v1';

type HealthResponse =
  paths['/healthz']['get']['responses']['200']['content']['application/json']['data'];

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

type ClientsResponse =
  paths['/clients']['get']['responses']['200']['content']['application/json']['data'];

export const fetchClients = async (tenantId: string): Promise<ClientsResponse> => {
  const res = await fetch(`${API_BASE.replace(/\/$/, '')}/clients`, {
    headers: {
      'X-Tenant-ID': tenantId,
    },
  });
  if (!res.ok) {
    throw new Error(`Falha ao consultar clientes: ${res.status}`);
  }
  const payload = (await res.json()) as {
    data: ClientsResponse;
  };
  return payload.data;
}

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
}

type ClientRequest = components['schemas']['ClientRequest'];
type ClientResponse = components['schemas']['Client'];

export const createClient = async (tenantId: string, client: ClientRequest): Promise<ClientResponse> => {
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
  const payload = (await res.json()) as {
    data: ClientResponse;
  };
  return payload.data;
}

export const updateClient = async (tenantId: string, clientId: string, client: ClientRequest): Promise<ClientResponse> => {
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
  const payload = (await res.json()) as {
    data: ClientResponse;
  };
  return payload.data;
}