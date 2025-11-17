import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';
import { useAuth } from '../../contexts/useAuth';
import { useSnackbar } from '../../hooks/useSnackbar';

// TODO: Substituir pelo tipo Product real e pelas funções da API
interface Product {
  id: string;
  name: string;
  price: number;
  stock: number;
}

interface ProductRequest {
  name: string;
  price: number;
  stock: number;
}

const createProduct = async (params: { tenantId: string; input: ProductRequest; accessToken: string }): Promise<Product> => {
  console.log('Creating product with params:', params);
  return { ...params.input, id: new Date().toISOString() };
};

const updateProduct = async (params: { tenantId: string; productId: string; input: ProductRequest; accessToken: string }): Promise<Product> => {
  console.log('Updating product with params:', params);
  return { ...params.input, id: params.productId };
};

interface ProductFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (product: Product) => void;
  product: Product | null;
}

const ProductForm: React.FC<ProductFormProps> = ({ open, onClose, onSave, product }) => {
  const [name, setName] = useState('');
  const [price, setPrice] = useState(0);
  const [stock, setStock] = useState(0);
  const { tenantId, accessToken } = useAuth();
  const { showSnackbar } = useSnackbar();

  useEffect(() => {
    if (product) {
      setName(product.name);
      setPrice(product.price);
      setStock(product.stock);
    } else {
      setName('');
      setPrice(0);
      setStock(0);
    }
  }, [product]);

  const handleSave = async () => {
    if (!tenantId || !accessToken) {
      return;
    }
    const productData: ProductRequest = { name, price, stock };
    try {
      let savedProduct;
      if (product) {
        savedProduct = await updateProduct({
          tenantId,
          productId: product.id,
          input: productData,
          accessToken,
        });
      } else {
        savedProduct = await createProduct({ tenantId, input: productData, accessToken });
      }
      onSave(savedProduct);
      onClose();
      showSnackbar(`Produto ${product ? 'atualizado' : 'criado'} com sucesso!`, 'success');
    } catch (error) {
      console.error('Error saving product:', error);
      showSnackbar('Erro ao salvar produto.', 'error');
    }
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{product ? 'Editar Produto' : 'Adicionar Produto'}</DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          id="name"
          label="Nome do Produto"
          type="text"
          fullWidth
          variant="standard"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <TextField
          margin="dense"
          id="price"
          label="Preço"
          type="number"
          fullWidth
          variant="standard"
          value={price}
          onChange={(e) => setPrice(Number(e.target.value))}
        />
        <TextField
          margin="dense"
          id="stock"
          label="Estoque"
          type="number"
          fullWidth
          variant="standard"
          value={stock}
          onChange={(e) => setStock(Number(e.target.value))}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancelar</Button>
        <Button onClick={handleSave}>Salvar</Button>
      </DialogActions>
    </Dialog>
  );
};

export default ProductForm;
