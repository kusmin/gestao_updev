import React, { useEffect } from 'react';
import { useForm, Controller, SubmitHandler } from 'react-hook-form';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';

import { Product } from '@/types/product';

interface ProductFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (product: Partial<Product>) => Promise<void>;
  product: Product | null;
}

const ProductForm: React.FC<ProductFormProps> = ({ open, onClose, onSave, product }) => {
  const { control, handleSubmit, reset } = useForm<Product>({
    defaultValues: {
      name: '',
      sku: '',
      price: 0,
      cost: 0,
      stock_qty: 0,
      min_stock: 0,
      description: '',
      tenant_id: '',
    },
  });

  useEffect(() => {
    if (open) {
      reset(
        product || {
          name: '',
          sku: '',
          price: 0,
          cost: 0,
          stock_qty: 0,
          min_stock: 0,
          description: '',
          tenant_id: '',
        },
      );
    }
  }, [product, open, reset]);

  const handleSave: SubmitHandler<Product> = (data) => {
    onSave(data).then(() => {
      onClose();
    });
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{product ? 'Editar Produto' : 'Adicionar Produto'}</DialogTitle>
      <form onSubmit={handleSubmit(handleSave)}>
        <DialogContent>
          <Controller
            name="name"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                autoFocus
                margin="dense"
                label="Nome"
                type="text"
                fullWidth
                variant="standard"
              />
            )}
          />
          <Controller
            name="sku"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                margin="dense"
                label="SKU"
                type="text"
                fullWidth
                variant="standard"
              />
            )}
          />
          <Controller
            name="price"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                margin="dense"
                label="Preço"
                type="number"
                fullWidth
                variant="standard"
              />
            )}
          />
          <Controller
            name="cost"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                margin="dense"
                label="Custo"
                type="number"
                fullWidth
                variant="standard"
              />
            )}
          />
          <Controller
            name="stock_qty"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                margin="dense"
                label="Estoque"
                type="number"
                fullWidth
                variant="standard"
              />
            )}
          />
          <Controller
            name="min_stock"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                margin="dense"
                label="Estoque Mínimo"
                type="number"
                fullWidth
                variant="standard"
              />
            )}
          />
          <Controller
            name="description"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                margin="dense"
                label="Descrição"
                type="text"
                fullWidth
                multiline
                rows={4}
                variant="standard"
              />
            )}
          />
          <Controller
            name="tenant_id"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                margin="dense"
                label="Tenant ID"
                type="text"
                fullWidth
                variant="standard"
              />
            )}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Cancelar</Button>
          <Button type="submit">Salvar</Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default ProductForm;
