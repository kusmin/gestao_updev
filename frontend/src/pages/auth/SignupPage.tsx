import { Box, Button, Link, Paper, Stack, TextField, Typography } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';

const SignupPage: React.FC = () => {
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
      <Paper elevation={3} sx={{ maxWidth: 480, width: '100%', p: 4 }}>
        <Stack spacing={3}>
          <div>
            <Typography variant="h5" component="h1" gutterBottom>
              Criar conta
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Informe os dados da empresa e do responsável para iniciar a plataforma.
            </Typography>
          </div>
          <TextField label="Nome da empresa" fullWidth />
          <TextField label="Documento" fullWidth />
          <TextField label="Telefone comercial" fullWidth />
          <TextField label="Nome do administrador" fullWidth />
          <TextField type="email" label="E-mail do administrador" fullWidth />
          <TextField type="password" label="Senha" fullWidth />
          <Button variant="contained" color="primary" size="large">
            Criar conta
          </Button>
          <Typography variant="body2" color="text.secondary" align="center">
            Já possui acesso?{' '}
            <Link component={RouterLink} to="/login">
              Fazer login
            </Link>
          </Typography>
        </Stack>
      </Paper>
    </Box>
  );
};

export default SignupPage;
