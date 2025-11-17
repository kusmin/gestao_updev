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
import { useAuth } from '../../contexts/useAuth';
import SaleForm from './SaleForm';

// TODO: Substituir pelo tipo Sale real e pelas funções da API
interface Sale {
  id: string;
  clientName: string;
  total: number;
  date: string;
}

const fetchSales = async (params: { tenantId: string; accessToken: string }): Promise<Sale[]> => {
  console.log('Fetching sales with params:', params);
  // Placeholder data
  return [
    { id: '1', clientName: 'João da Silva', total: 50.0, date: '2025-11-18T10:30:00Z' },
    { id: '2', clientName: 'Maria Oliveira', total: 35.0, date: '2025-11-18T11:30:00Z' },
  ];
};

const deleteSale = async (params: { tenantId: string; saleId: string; accessToken: string }): Promise<void> => {
  console.log('Deleting sale with params:', params);
  return Promise.resolve();
};

const SaleListPage: React.FC = () => {
  const [sales, setSales] = useState<Sale[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingSale, setEditingSale] = useState<Sale | null>(null);
  const { tenantId, accessToken } = useAuth();

  const getSales = useCallback(async () => {
    if (!tenantId || !accessToken) {
      setSales([]);
      return;
    }
    try {
      const data = await fetchSales({ tenantId, accessToken });
      if (data && Array.isArray(data)) {
        setSales(data);
      }
    } catch (error) {
      console.error('Error fetching sales:', error);
    }
  }, [tenantId, accessToken, setSales]);

  useEffect(() => {
    if (tenantId && accessToken) {
      getSales();
    }
  }, [tenantId, accessToken, getSales]);

  const handleOpenForm = (sale: Sale | null = null) => {
    setEditingSale(sale);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingSale(null);
    setIsFormOpen(false);
  };

  const handleSaveSale = (sale: Sale) => {
    if (editingSale) {
      setSales(sales.map((s) => (s.id === sale.id ? sale : s)));
    } else {
      setSales([...sales, sale]);
    }
    handleCloseForm();
  };

  const handleDelete = async (id: string) => {
    if (!tenantId || !accessToken) {
      return;
    }
    try {
      await deleteSale({ tenantId, saleId: id, accessToken });
      setSales((prev) => prev.filter((sale) => sale.id !== id));
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
          Registrar Venda
        </Button>
      </Box>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Cliente</TableCell>
              <TableCell>Data</TableCell>
              <TableCell>Total</TableCell>
              <TableCell>Ações</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {sales.map((sale) => (
              <TableRow key={sale.id}>
                <TableCell>{sale.clientName}</TableCell>
                <TableCell>{new Date(sale.date).toLocaleString()}</TableCell>
                <TableCell>R$ {sale.total.toFixed(2)}</TableCell>
                <TableCell>
                  <Button onClick={() => handleOpenForm(sale)}>Detalhes</Button>
                  <Button color="error" onClick={() => handleDelete(sale.id)}>
                    Cancelar
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
