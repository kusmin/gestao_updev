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
import { Link } from 'react-router-dom';

// Mock data for sales
const mockSales = [
  { id: '1', client: 'Cliente 1', total: 100.0, created_at: new Date().toISOString() },
  { id: '2', client: 'Cliente 2', total: 200.0, created_at: new Date().toISOString() },
  { id: '3', client: 'Cliente 3', total: 300.0, created_at: new Date().toISOString() },
];

const SaleListPage: React.FC = () => {
  const [isLoading] = useState(false);
  const [isError] = useState(false);

  return (
    <Container>
      <Box sx={{ my: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" component="h1">
          Vendas
        </Typography>
        <Button variant="contained" color="primary" component={Link} to="/sales/new">
          Adicionar Venda
        </Button>
      </Box>

      {isLoading && <CircularProgress />}
      {isError && <Alert severity="error">Erro ao carregar vendas.</Alert>}

      {!isLoading && !isError && (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Cliente</TableCell>
                <TableCell>Total</TableCell>
                <TableCell>Data</TableCell>
                <TableCell>Ações</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {mockSales.map((sale) => (
                <TableRow key={sale.id}>
                  <TableCell>{sale.client}</TableCell>
                  <TableCell>{sale.total.toFixed(2)}</TableCell>
                  <TableCell>{new Date(sale.created_at).toLocaleDateString()}</TableCell>
                  <TableCell>
                    <Button component={Link} to={`/sales/edit/${sale.id}`}>
                      Editar
                    </Button>
                    <Button color="error" onClick={() => alert('Excluir ' + sale.id)}>
                      Excluir
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}
    </Container>
  );
};

export default SaleListPage;