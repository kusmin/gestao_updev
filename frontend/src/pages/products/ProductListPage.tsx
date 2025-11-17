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
import ProductForm from './ProductForm';

// TODO: Substituir pelo tipo Product real e pelas funções da API
interface Product {
  id: string;
  name: string;
  price: number;
  stock: number;
}

const fetchProducts = async (params: { tenantId: string; accessToken: string }): Promise<Product[]> => {
  console.log('Fetching products with params:', params);
  // Placeholder data
  return [
    { id: '1', name: 'Pomada Modeladora', price: 25.0, stock: 50 },
    { id: '2', name: 'Óleo para Barba', price: 35.0, stock: 30 },
  ];
};

const deleteProduct = async (params: { tenantId: string; productId: string; accessToken: string }): Promise<void> => {
  console.log('Deleting product with params:', params);
  return Promise.resolve();
};

const ProductListPage: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const { tenantId, accessToken } = useAuth();

  const getProducts = useCallback(async () => {
    if (!tenantId || !accessToken) {
      setProducts([]);
      return;
    }
    try {
      const data = await fetchProducts({ tenantId, accessToken });
      if (data && Array.isArray(data)) {
        setProducts(data);
      }
    } catch (error) {
      console.error('Error fetching products:', error);
    }
  }, [tenantId, accessToken, setProducts]);

  useEffect(() => {
    if (tenantId && accessToken) {
      getProducts();
    }
  }, [tenantId, accessToken, getProducts]);

  const handleOpenForm = (product: Product | null = null) => {
    setEditingProduct(product);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingProduct(null);
    setIsFormOpen(false);
  };

  const handleSaveProduct = (product: Product) => {
    if (editingProduct) {
      setProducts(products.map((p) => (p.id === product.id ? product : p)));
    } else {
      setProducts([...products, product]);
    }
    handleCloseForm();
  };

  const handleDelete = async (id: string) => {
    if (!tenantId || !accessToken) {
      return;
    }
    try {
      await deleteProduct({ tenantId, productId: id, accessToken });
      setProducts((prev) => prev.filter((product) => product.id !== id));
    } catch (error) {
      console.error('Error deleting product:', error);
    }
  };

  return (
    <Container>
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Produtos
        </Typography>
        <Button variant="contained" color="primary" onClick={() => handleOpenForm()}>
          Adicionar Produto
        </Button>
      </Box>
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
            {products.map((product) => (
              <TableRow key={product.id}>
                <TableCell>{product.name}</TableCell>
                <TableCell>R$ {product.price.toFixed(2)}</TableCell>
                <TableCell>{product.stock}</TableCell>
                <TableCell>
                  <Button onClick={() => handleOpenForm(product)}>Editar</Button>
                  <Button color="error" onClick={() => handleDelete(product.id)}>
                    Excluir
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <ProductForm
        open={isFormOpen}
        onClose={handleCloseForm}
        onSave={handleSaveProduct}
        product={editingProduct}
      />
    </Container>
  );
};

export default ProductListPage;
