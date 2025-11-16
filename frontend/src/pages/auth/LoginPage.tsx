import { Box, Button, Link, Paper, Stack, TextField, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

const LoginPage: React.FC = () => {
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
        <Stack spacing={3}>
          <div>
            <Typography variant="h5" component="h1" gutterBottom>
              Entrar
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Acesse o painel para gerenciar sua operação.
            </Typography>
          </div>
          <TextField type="email" label="E-mail" fullWidth />
          <TextField type="password" label="Senha" fullWidth />
          <Button variant="contained" color="primary" size="large">
            Entrar
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
