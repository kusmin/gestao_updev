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
import SaleForm from './SaleForm';
import apiClient from '../../lib/apiClient';

interface Sale {
  id: string;
  client_id: string;
  total: number;
  status: string;
  tenant_id: string;
}

const SaleListPage: React.FC = () => {
  const [sales, setSales] = useState<Sale[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingSale, setEditingSale] = useState<Sale | null>(null);

  const fetchSales = async () => {
    try {
      const response = await apiClient<{ data: Sale[] }>('/admin/sales/orders');
      setSales(response.data);
    } catch (error) {
      console.error('Error fetching sales:', error);
    }
  };

  useEffect(() => {
    fetchSales();
  }, []);

  const handleOpenForm = (sale: Sale | null = null) => {
    setEditingSale(sale);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingSale(null);
    setIsFormOpen(false);
    fetchSales();
  };

  const handleSaveSale = async (sale: Partial<Sale>) => {
    try {
      if (editingSale) {
        await apiClient(`/admin/sales/orders/${editingSale.id}`, {
          method: 'PUT',
          body: JSON.stringify(sale),
        });
      } else {
        await apiClient('/admin/sales/orders', {
          method: 'POST',
          body: JSON.stringify(sale),
        });
      }
    } catch (error) {
      console.error('Error saving sale:', error);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiClient(`/admin/sales/orders/${id}`, { method: 'DELETE' });
      fetchSales();
    } catch (error) {
      console.error('Error deleting sale:', error);
    }
  };

  return (
    <Container>
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Vendas
        </Typography>
        <Button variant="contained" color="primary" onClick={() => handleOpenForm()}>
          Adicionar Venda
        </Button>
      </Box>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Cliente ID</TableCell>
              <TableCell>Total</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Tenant ID</TableCell>
              <TableCell>Ações</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {sales.map((sale) => (
              <TableRow key={sale.id}>
                <TableCell>{sale.client_id}</TableCell>
                <TableCell>{sale.total}</TableCell>
                <TableCell>{sale.status}</TableCell>
                <TableCell>{sale.tenant_id}</TableCell>
                <TableCell>
                  <Button onClick={() => handleOpenForm(sale)}>Editar</Button>
                  <Button color="error" onClick={() => handleDelete(sale.id)}>
                    Excluir
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <SaleForm
        open={isFormOpen}
        onClose={handleCloseForm}
        onSave={handleSaveSale}
        sale={editingSale}
      />
    </Container>
  );
};

export default SaleListPage;