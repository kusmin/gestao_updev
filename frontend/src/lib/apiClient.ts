import type { paths } from '../types/api';

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