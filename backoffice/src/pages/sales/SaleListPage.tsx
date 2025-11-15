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
import apiClient from '@/lib/apiClient';
import { SalesOrder } from '@/types/sales';

const SaleListPage: React.FC = () => {
  const [sales, setSales] = useState<SalesOrder[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingSale, setEditingSale] = useState<SalesOrder | null>(null);

  const fetchSales = async () => {
    try {
      const response = await apiClient<{ data: SalesOrder[] }>('/admin/sales/orders');
      setSales(response.data);
    } catch (error) {
      console.error('Error fetching sales:', error);
    }
  };

  useEffect(() => {
    fetchSales();
  }, []);

  const handleOpenForm = (sale: SalesOrder | null = null) => {
    setEditingSale(sale);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingSale(null);
    setIsFormOpen(false);
    fetchSales(); // Refetch sales after closing form
  };

  const handleSaveSale = async (sale: Partial<SalesOrder>) => {
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
      fetchSales(); // Refetch sales after deleting
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
              <TableCell>Cliente</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Total</TableCell>
              <TableCell>Ações</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {sales.map((sale) => (
              <TableRow key={sale.id}>
                <TableCell>{sale.client_id}</TableCell>
                <TableCell>{sale.status}</TableCell>
                <TableCell>{sale.total}</TableCell>
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
