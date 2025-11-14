const BASE_URL = 'http://localhost:8080/v1'; // TODO: use env var

let authToken: string | null = null;

export const setAuthToken = (token: string | null) => {
  authToken = token;
};

interface ApiClientOptions extends RequestInit {
  tenantId?: string;
}

const apiClient = async <T>(
  endpoint: string,
  options: ApiClientOptions = {}
): Promise<T> => {
  const { headers, tenantId, ...rest } = options;

  const defaultHeaders: Record<string, string> = {
    'Content-Type': 'application/json',
    ...headers,
  };

  if (authToken) {
    defaultHeaders['Authorization'] = `Bearer ${authToken}`;
  }

  if (tenantId) {
    defaultHeaders['X-Tenant-ID'] = tenantId;
  }

  const response = await fetch(`${BASE_URL}${endpoint}`, {
    headers: defaultHeaders,
    ...rest,
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.message || 'Something went wrong');
  }

  if (response.status === 204) {
    return {} as T;
  }

  return response.json();
};

export default apiClient;