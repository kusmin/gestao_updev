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

interface Client {
  id: string;
  name: string;
  email: string;
  phone: string;
  notes: string;
  tags: string[];
  contact: Record<string, any>;
}

const ClientListPage: React.FC = () => {
  const [clients, setClients] = useState<Client[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingClient, setEditingClient] = useState<Client | null>(null);

  const fetchClients = async () => {
    try {
      // TODO: Replace with actual API endpoint and add authentication
      const response = await fetch('http://localhost:8080/v1/clients', {
        headers: {
          // TODO: Replace with actual tenant ID
          'X-Tenant-ID': 'a4b2b2b2-b2b2-4b2b-b2b2-b2b2b2b2b2b2',
        },
      });
      const data = await response.json();
      setClients(data.data); // Assuming the API returns data in a 'data' property
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
  };

  const handleSaveClient = (client: Client) => {
    if (editingClient) {
      setClients(
        clients.map((c) => (c.id === client.id ? client : c))
      );
    } else {
      setClients([...clients, client]);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      // TODO: Replace with actual API endpoint and add authentication
      await fetch(`http://localhost:8080/v1/clients/${id}`, {
        method: 'DELETE',
        headers: {
          // TODO: Replace with actual tenant ID
          'X-Tenant-ID': 'a4b2b2b2-b2b2-4b2b-b2b2-b2b2b2b2b2b2',
        },
      });
      setClients(clients.filter((client) => client.id !== id));
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