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

// Mock data for products
const mockProducts = [
  { id: '1', name: 'Produto 1', price: 10.0, stock_qty: 100 },
  { id: '2', name: 'Produto 2', price: 20.0, stock_qty: 200 },
  { id: '3', name: 'Produto 3', price: 30.0, stock_qty: 300 },
];

const ProductListPage: React.FC = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);

  return (
    <Container>
      <Box sx={{ my: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" component="h1">
          Produtos
        </Typography>
        <Button variant="contained" color="primary" component={Link} to="/products/new">
          Adicionar Produto
        </Button>
      </Box>

      {isLoading && <CircularProgress />}
      {isError && <Alert severity="error">Erro ao carregar produtos.</Alert>}

      {!isLoading && !isError && (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Nome</TableCell>
                <TableCell>Preço</TableCell>
                <TableCell>Estoque</TableCell>
                <TableCell>Ações</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {mockProducts.map((product) => (
                <TableRow key={product.id}>
                  <TableCell>{product.name}</TableCell>
                  <TableCell>{product.price.toFixed(2)}</TableCell>
                  <TableCell>{product.stock_qty}</TableCell>
                  <TableCell>
                    <Button component={Link} to={`/products/edit/${product.id}`}>
                      Editar
                    </Button>
                    <Button color="error" onClick={() => alert('Excluir ' + product.name)}>
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

export default ProductListPage;