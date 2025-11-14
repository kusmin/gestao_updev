import React, { useEffect, useState } from 'react';
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
import ClientForm from './ClientForm';
import apiClient from '../../lib/apiClient';

interface Client {
  id: string;
  name: string;
  email: string;
  phone: string;
  notes: string;
  tags: string[];
  contact: Record<string, unknown>;
  tenant_id: string;
}

const ClientListPage: React.FC = () => {
  const [clients, setClients] = useState<Client[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingClient, setEditingClient] = useState<Client | null>(null);

  const fetchClients = async () => {
    try {
      const response = await apiClient<{ data: Client[] }>('/admin/clients');
      setClients(response.data);
    } catch (error) {
      console.error('Error fetching clients:', error);
    }
  };

  useEffect(() => {
    fetchClients();
  }, []);

  const handleOpenForm = (client: Client | null = null) => {
    setEditingClient(client);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingClient(null);
    setIsFormOpen(false);
    fetchClients();
  };

  const handleSaveClient = async (client: Partial<Client>) => {
    try {
      if (editingClient) {
        await apiClient(`/admin/clients/${editingClient.id}`, {
          method: 'PUT',
          body: JSON.stringify(client),
        });
      } else {
        await apiClient('/admin/clients', {
          method: 'POST',
          body: JSON.stringify(client),
        });
      }
    } catch (error) {
      console.error('Error saving client:', error);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiClient(`/admin/clients/${id}`, { method: 'DELETE' });
      fetchClients();
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
              <TableCell>Tenant ID</TableCell>
              <TableCell>Ações</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {clients.map((client) => (
              <TableRow key={client.id}>
                <TableCell>{client.name}</TableCell>
                <TableCell>{client.email}</TableCell>
                <TableCell>{client.phone}</TableCell>
                <TableCell>{client.tenant_id}</TableCell>
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
