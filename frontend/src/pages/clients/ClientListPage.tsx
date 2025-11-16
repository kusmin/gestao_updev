import React, { useCallback, useEffect, useState } from 'react';
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
} from '@mui/material';
import { deleteClient, fetchClients, type Client } from '../../lib/apiClient';
import ClientForm from './ClientForm';
import { useAuth } from '../../contexts/useAuth';

const ClientListPage: React.FC = () => {
  const [clients, setClients] = useState<Client[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingClient, setEditingClient] = useState<Client | null>(null);
  const { tenantId, accessToken } = useAuth();

  const getClients = useCallback(async () => {
    if (!tenantId || !accessToken) {
      setClients([]);
      return;
    }
    try {
      const data = await fetchClients({ tenantId, accessToken });
      if (data && Array.isArray(data)) {
        setClients(data);
      }
    } catch (error) {
      console.error('Error fetching clients:', error);
    }
  }, [tenantId, accessToken, setClients]);

  useEffect(() => {
    if (tenantId && accessToken) {
      getClients();
    }
  }, [tenantId, accessToken, getClients]);

  const handleOpenForm = (client: Client | null = null) => {
    setEditingClient(client);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingClient(null);
    setIsFormOpen(false);
  };

  const handleSaveClient = (client: Client) => {
    if (editingClient) {
      setClients(clients.map((c) => (c.id === client.id ? client : c)));
    } else {
      setClients([...clients, client]);
    }
  };

  const handleDelete = async (id: string) => {
    if (!tenantId || !accessToken) {
      return;
    }
    try {
      await deleteClient({ tenantId, clientId: id, accessToken });
      setClients((prev) => prev.filter((client) => client.id !== id));
    } catch (error) {
      console.error('Error deleting client:', error);
    }
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
            {clients.map((client) => (
              <TableRow key={client.id}>
                <TableCell>{client.name}</TableCell>
                <TableCell>{client.email}</TableCell>
                <TableCell>{client.phone}</TableCell>
                <TableCell>
                  <Button onClick={() => handleOpenForm(client)}>Editar</Button>
                  <Button color="error" onClick={() => handleDelete(client.id)}>
                    Excluir
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <ClientForm
        open={isFormOpen}
        onClose={handleCloseForm}
        onSave={handleSaveClient}
        client={editingClient}
      />
    </Container>
  );
};

export default ClientListPage;
