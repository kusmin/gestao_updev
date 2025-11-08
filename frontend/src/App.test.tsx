import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { CssBaseline, ThemeProvider } from '@mui/material';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';

import App from './App';
import theme from './lib/theme';

const mockFetchHealth = vi.fn();

vi.mock('./lib/apiClient', () => ({
  fetchHealth: (...args: unknown[]) => mockFetchHealth(...args),
}));

const renderWithProviders = () => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <App />
      </ThemeProvider>
    </QueryClientProvider>,
  );
};

describe('App', () => {
  beforeEach(() => {
    mockFetchHealth.mockReset();
  });

  it('exibe os dados retornados pela API quando a chamada é bem-sucedida', async () => {
    mockFetchHealth.mockResolvedValue({
      status: 'OK',
      env: 'local',
    });

    renderWithProviders();

    expect(
      await screen.findByRole('heading', { name: /gestão updev/i }),
    ).toBeInTheDocument();
    expect(await screen.findByText(/ambiente: local/i)).toBeInTheDocument();
  });

  it('mostra alerta de erro quando a chamada falha', async () => {
    mockFetchHealth.mockRejectedValueOnce(new Error('offline'));

    renderWithProviders();

    expect(
      await screen.findByText(/não foi possível conectar à api/i),
    ).toBeInTheDocument();
  });

  it('permite refazer a consulta ao focar novamente (exemplo com userEvent)', async () => {
    const user = userEvent.setup();
    mockFetchHealth.mockResolvedValue({ status: 'OK', env: 'local' });

    renderWithProviders();

    const envInfo = await screen.findByText(/ambiente: local/i);
    await user.tab();

    expect(envInfo).toBeInTheDocument();
    expect(mockFetchHealth).toHaveBeenCalledTimes(1);
  });
});
