import {
  Alert,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Container,
  Stack,
  Typography,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';

import { fetchHealth } from './lib/apiClient';

function App() {
  const { data, isLoading, isError } = useQuery({
    queryKey: ['health'],
    queryFn: fetchHealth,
  });

  const statusChip = data
    ? (() => {
        const normalized = data.status?.toLowerCase();
        const color =
          normalized === 'ok'
            ? 'success'
            : normalized === 'degraded'
              ? 'warning'
              : 'error';

        return (
          <Chip
            label={data.status}
            color={color}
            variant="filled"
            sx={{ textTransform: 'uppercase', letterSpacing: 0.5 }}
          />
        );
      })()
    : null;

  return (
    <Container maxWidth="sm" sx={{ py: 8 }}>
      <Card elevation={8} sx={{ borderRadius: 3 }}>
        <CardContent>
          <Stack spacing={3} textAlign="center">
            <Stack spacing={1}>
              <Typography variant="h3" component="h1" fontWeight={600}>
                Gestão UpDev
              </Typography>
              <Typography color="text.secondary">
                Plataforma de gestão para negócios locais.
              </Typography>
            </Stack>

            {isLoading && (
              <Stack alignItems="center" spacing={1}>
                <CircularProgress size={32} />
                <Typography variant="body2" color="text.secondary">
                  Carregando status da API...
                </Typography>
              </Stack>
            )}

            {isError && (
              <Alert severity="error" variant="outlined">
                Não foi possível conectar à API. Verifique se o backend está
                rodando.
              </Alert>
            )}

            {!isLoading && !isError && data && (
              <Stack
                direction={{ xs: 'column', sm: 'row' }}
                spacing={2}
                justifyContent="center"
                alignItems="center"
              >
                <Typography variant="subtitle1" fontWeight={500}>
                  API
                </Typography>
                {statusChip}
                <Typography variant="body2" color="text.secondary">
                  Ambiente: {data.env}
                </Typography>
              </Stack>
            )}
          </Stack>
        </CardContent>
      </Card>
    </Container>
  );
}

export default App;
