import React, { useState } from 'react';
import {
  Box,
  Button,
  Container,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  CircularProgress,
  Alert,
} from '@mui/material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { deleteClient, fetchClients, type Client } from '../../lib/apiClient';
import ClientForm from './ClientForm';
import { useAuth } from '../../contexts/useAuth';
import { useSnackbar } from '../../hooks/useSnackbar';

const ClientListPage: React.FC = () => {
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingClient, setEditingClient] = useState<Client | null>(null);
  const { tenantId, accessToken } = useAuth();
  const { showSnackbar } = useSnackbar();
  const queryClient = useQueryClient();

  const { data: clients, isLoading, isError } = useQuery({
    queryKey: ['clients', tenantId],
    queryFn: () => fetchClients({ tenantId: tenantId!, accessToken: accessToken! }),
    enabled: !!tenantId && !!accessToken,
  });

  const deleteMutation = useMutation({
    mutationFn: (clientId: string) => deleteClient({ tenantId: tenantId!, clientId, accessToken: accessToken! }),
    onSuccess: () => {
      showSnackbar('Cliente excluído com sucesso!', 'success');
      queryClient.invalidateQueries({ queryKey: ['clients', tenantId] });
    },
    onError: () => {
      showSnackbar('Erro ao excluir cliente.', 'error');
    },
  });

  const handleOpenForm = (client: Client | null = null) => {
    setEditingClient(client);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingClient(null);
    setIsFormOpen(false);
  };

  const handleSaveSuccess = () => {
    queryClient.invalidateQueries({ queryKey: ['clients', tenantId] });
  };

  const handleDelete = (id: string) => {
    deleteMutation.mutate(id);
  };

  return (
    <Container>
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Clientes
        </Typography>
        <Button variant="contained" color="primary" onClick={() => handleOpenForm()}>
          Adicionar Cliente
        </Button>
      </Box>

      {isLoading && <CircularProgress />}
      {isError && <Alert severity="error">Erro ao carregar clientes.</Alert>}

      {!isLoading && !isError && (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Nome</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Telefone</TableCell>
                <TableCell>Ações</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {clients?.map((client) => (
                <TableRow key={client.id}>
                  <TableCell>{client.name}</TableCell>
                  <TableCell>{client.email}</TableCell>
                  <TableCell>{client.phone}</TableCell>
                  <TableCell>
                    <Button onClick={() => handleOpenForm(client)}>Editar</Button>
                    <Button color="error" onClick={() => handleDelete(client.id)} disabled={deleteMutation.isPending}>
                      Excluir
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}

      <ClientForm
        open={isFormOpen}
        onClose={handleCloseForm}
        onSave={handleSaveSuccess}
        client={editingClient}
      />
    </Container>
  );
};

export default ClientListPage;
