import { Box, Button, Link, Paper, Stack, TextField, Typography, Alert } from '@mui/material';
import { Link as RouterLink, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { useEffect, useState } from 'react';

const LoginPage: React.FC = () => {
  const { login, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const redirectTo = (location.state as { from?: string })?.from ?? '/';
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/', { replace: true });
    }
  }, [isAuthenticated, navigate]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      await login({ email, password });
      navigate(redirectTo, { replace: true });
    } catch (err) {
      console.error(err);
      setError('Credenciais inválidas, tente novamente.');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: (theme) => theme.palette.background.default,
        p: 2,
      }}
    >
      <Paper elevation={3} sx={{ maxWidth: 420, width: '100%', p: 4 }}>
        <Stack spacing={3} component="form" onSubmit={handleSubmit}>
          <div>
            <Typography variant="h5" component="h1" gutterBottom>
              Entrar
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Acesse o painel para gerenciar sua operação.
            </Typography>
          </div>
          {error && <Alert severity="error">{error}</Alert>}
          <TextField
            type="email"
            label="E-mail"
            fullWidth
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            required
          />
          <TextField
            type="password"
            label="Senha"
            fullWidth
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            required
          />
          <Button variant="contained" color="primary" size="large" type="submit" disabled={submitting}>
            {submitting ? 'Entrando...' : 'Entrar'}
          </Button>
          <Typography variant="body2" color="text.secondary" align="center">
            Ainda não tem conta?{' '}
            <Link component={RouterLink} to="/signup">
              Criar conta
            </Link>
          </Typography>
        </Stack>
      </Paper>
    </Box>
  );
};

export default LoginPage;
