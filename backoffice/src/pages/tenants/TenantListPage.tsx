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
import TenantForm from './TenantForm';
import apiClient from '../../lib/apiClient';

export interface Tenant {
  id: string;
  name: string;
  document: string;
  email: string;
  phone: string;
}

const TenantListPage: React.FC = () => {
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingTenant, setEditingTenant] = useState<Tenant | null>(null);

  const fetchTenants = async () => {
    try {
      const response = await apiClient<{ data: Tenant[] }>('/admin/tenants');
      setTenants(response.data);
    } catch (error) {
      console.error('Error fetching tenants:', error);
    }
  };

  useEffect(() => {
    fetchTenants();
  }, []);

  const handleOpenForm = (tenant: Tenant | null = null) => {
    setEditingTenant(tenant);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingTenant(null);
    setIsFormOpen(false);
    fetchTenants();
  };

  const handleSaveTenant = async (tenant: Partial<Tenant>) => {
    try {
      if (editingTenant) {
        await apiClient(`/admin/tenants/${editingTenant.id}`, {
          method: 'PUT',
          body: JSON.stringify(tenant),
        });
      } else {
        await apiClient('/admin/tenants', {
          method: 'POST',
          body: JSON.stringify(tenant),
        });
      }
    } catch (error) {
      console.error('Error saving tenant:', error);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiClient(`/admin/tenants/${id}`, { method: 'DELETE' });
      fetchTenants();
    } catch (error) {
      console.error('Error deleting tenant:', error);
    }
  };

  return (
    <Container>
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Tenants
        </Typography>
        <Button variant="contained" color="primary" onClick={() => handleOpenForm()}>
          Adicionar Tenant
        </Button>
      </Box>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Nome</TableCell>
              <TableCell>Documento</TableCell>
              <TableCell>Email</TableCell>
              <TableCell>Telefone</TableCell>
              <TableCell>Ações</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {tenants.map((tenant) => (
              <TableRow key={tenant.id}>
                <TableCell>{tenant.name}</TableCell>
                <TableCell>{tenant.document}</TableCell>
                <TableCell>{tenant.email}</TableCell>
                <TableCell>{tenant.phone}</TableCell>
                <TableCell>
                  <Button onClick={() => handleOpenForm(tenant)}>Editar</Button>
                  <Button color="error" onClick={() => handleDelete(tenant.id)}>
                    Excluir
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <TenantForm
        open={isFormOpen}
        onClose={handleCloseForm}
        onSave={handleSaveTenant}
        tenant={editingTenant}
      />
    </Container>
  );
};

export default TenantListPage;