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
import ProductForm from './ProductForm';
import apiClient from '../../lib/apiClient';

interface Product {
  id: string;
  name: string;
  sku: string;
  price: number;
  stock_qty: number;
  tenant_id: string;
}

const ProductListPage: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);

  const fetchProducts = async () => {
    try {
      const response = await apiClient<{ data: Product[] }>('/admin/products');
      setProducts(response.data);
    } catch (error) {
      console.error('Error fetching products:', error);
    }
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  const handleOpenForm = (product: Product | null = null) => {
    setEditingProduct(product);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingProduct(null);
    setIsFormOpen(false);
    fetchProducts(); // Refetch products after closing form
  };

  const handleSaveProduct = async (product: Partial<Product>) => {
    try {
      if (editingProduct) {
        await apiClient(`/admin/products/${editingProduct.id}`, {
          method: 'PUT',
          body: JSON.stringify(product),
        });
      } else {
        await apiClient('/admin/products', {
          method: 'POST',
          body: JSON.stringify(product),
        });
      }
    } catch (error) {
      console.error('Error saving product:', error);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiClient(`/admin/products/${id}`, { method: 'DELETE' });
      fetchProducts(); // Refetch products after deleting
    } catch (error) {
      console.error('Error deleting product:', error);
    }
  };

  return (
    <Container>
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Products
        </Typography>
        <Button
          variant="contained"
          color="primary"
          onClick={() => handleOpenForm()}
        >
          Adicionar Produto
        </Button>
      </Box>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Nome</TableCell>
              <TableCell>SKU</TableCell>
              <TableCell>Preço</TableCell>
              <TableCell>Estoque</TableCell>
              <TableCell>Ações</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {products.map((product) => (
              <TableRow key={product.id}>
                <TableCell>{product.name}</TableCell>
                <TableCell>{product.sku}</TableCell>
                <TableCell>{product.price}</TableCell>
                <TableCell>{product.stock_qty}</TableCell>
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