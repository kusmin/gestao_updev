import { Box, Button, Link, Paper, Stack, TextField, Typography, Alert } from '@mui/material';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/useAuth';
import { useEffect, useState } from 'react';

const SignupPage: React.FC = () => {
  const { signup, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [companyName, setCompanyName] = useState('');
  const [companyDocument, setCompanyDocument] = useState('');
  const [companyPhone, setCompanyPhone] = useState('');
  const [userName, setUserName] = useState('');
  const [userEmail, setUserEmail] = useState('');
  const [userPhone, setUserPhone] = useState('');
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
      await signup({
        company: {
          name: companyName,
          document: companyDocument,
          phone: companyPhone,
        },
        user: {
          name: userName,
          email: userEmail,
          phone: userPhone,
          password,
        },
      });
      navigate('/', { replace: true });
    } catch (err) {
      console.error(err);
      setError('Não foi possível concluir o cadastro. Verifique os dados e tente novamente.');
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
      <Paper elevation={3} sx={{ maxWidth: 520, width: '100%', p: 4 }}>
        <Stack spacing={3} component="form" onSubmit={handleSubmit}>
          <div>
            <Typography variant="h5" component="h1" gutterBottom>
              Criar conta
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Informe os dados da empresa e do responsável para iniciar a plataforma.
            </Typography>
          </div>
          {error && <Alert severity="error">{error}</Alert>}
          <TextField
            label="Nome da empresa"
            fullWidth
            value={companyName}
            onChange={(event) => setCompanyName(event.target.value)}
            required
          />
          <TextField
            label="Documento"
            fullWidth
            value={companyDocument}
            onChange={(event) => setCompanyDocument(event.target.value)}
          />
          <TextField
            label="Telefone comercial"
            fullWidth
            value={companyPhone}
            onChange={(event) => setCompanyPhone(event.target.value)}
          />
          <TextField
            label="Nome do administrador"
            fullWidth
            value={userName}
            onChange={(event) => setUserName(event.target.value)}
            required
          />
          <TextField
            type="email"
            label="E-mail do administrador"
            fullWidth
            value={userEmail}
            onChange={(event) => setUserEmail(event.target.value)}
            required
          />
          <TextField
            label="Telefone do administrador"
            fullWidth
            value={userPhone}
            onChange={(event) => setUserPhone(event.target.value)}
          />
          <TextField
            type="password"
            label="Senha"
            fullWidth
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            required
          />
          <Button
            variant="contained"
            color="primary"
            size="large"
            type="submit"
            disabled={submitting}
          >
            {submitting ? 'Processando...' : 'Criar conta'}
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
