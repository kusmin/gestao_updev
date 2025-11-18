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
import ServiceForm from './ServiceForm';
import apiClient from '../../lib/apiClient';

interface Service {
  id: string;
  name: string;
  price: number;
  duration_minutes: number;
  tenant_id: string;
}

const ServiceListPage: React.FC = () => {
  const [services, setServices] = useState<Service[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingService, setEditingService] = useState<Service | null>(null);

  const fetchServices = async () => {
    try {
      const response = await apiClient<{ data: Service[] }>('/admin/services');
      setServices(response.data);
    } catch (error) {
      console.error('Error fetching services:', error);
    }
  };

  useEffect(() => {
    fetchServices();
  }, []);

  const handleOpenForm = (service: Service | null = null) => {
    setEditingService(service);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingService(null);
    setIsFormOpen(false);
    fetchServices();
  };

  const handleSaveService = async (service: Partial<Service>) => {
    try {
      if (editingService) {
        await apiClient(`/admin/services/${editingService.id}`, {
          method: 'PUT',
          body: JSON.stringify(service),
        });
      } else {
        await apiClient('/admin/services', {
          method: 'POST',
          body: JSON.stringify(service),
        });
      }
    } catch (error) {
      console.error('Error saving service:', error);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiClient(`/admin/services/${id}`, { method: 'DELETE' });
      fetchServices();
    } catch (error) {
      console.error('Error deleting service:', error);
    }
  };

  return (
    <Container>
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Serviços
        </Typography>
        <Button variant="contained" color="primary" onClick={() => handleOpenForm()}>
          Adicionar Serviço
        </Button>
      </Box>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Nome</TableCell>
              <TableCell>Preço</TableCell>
              <TableCell>Duração (min)</TableCell>
              <TableCell>Tenant ID</TableCell>
              <TableCell>Ações</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {services.map((service) => (
              <TableRow key={service.id}>
                <TableCell>{service.name}</TableCell>
                <TableCell>{service.price}</TableCell>
                <TableCell>{service.duration_minutes}</TableCell>
                <TableCell>{service.tenant_id}</TableCell>
                <TableCell>
                  <Button onClick={() => handleOpenForm(service)}>Editar</Button>
                  <Button color="error" onClick={() => handleDelete(service.id)}>
                    Excluir
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <ServiceForm
        open={isFormOpen}
        onClose={handleCloseForm}
        onSave={handleSaveService}
        service={editingService}
      />
    </Container>
  );
};

export default ServiceListPage;