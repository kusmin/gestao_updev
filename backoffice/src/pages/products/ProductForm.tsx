import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';

interface Product {
  id: string;
  name: string;
  sku: string;
  price: number;
  cost: number;
  stock_qty: number;
  min_stock: number;
  description: string;
  tenant_id: string;
}

interface ProductFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (product: Partial<Product>) => Promise<void>;
  product: Product | null;
}

const ProductForm: React.FC<ProductFormProps> = ({
  open,
  onClose,
  onSave,
  product,
}) => {
  const [formData, setFormData] = useState<Partial<Product>>({
    name: '',
    sku: '',
    price: 0,
    cost: 0,
    stock_qty: 0,
    min_stock: 0,
    description: '',
    tenant_id: '',
  });

  useEffect(() => {
    if (product) {
      setFormData(product);
    } else {
      setFormData({
        name: '',
        sku: '',
        price: 0,
        cost: 0,
        stock_qty: 0,
        min_stock: 0,
        description: '',
        tenant_id: '',
      });
    }
  }, [product, open]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    // TODO: Add validation
    await onSave(formData);
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{product ? 'Editar Produto' : 'Adicionar Produto'}</DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          name="name"
          label="Nome"
          type="text"
          fullWidth
          variant="standard"
          value={formData.name}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="sku"
          label="SKU"
          type="text"
          fullWidth
          variant="standard"
          value={formData.sku}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="price"
          label="Preço"
          type="number"
          fullWidth
          variant="standard"
          value={formData.price}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="cost"
          label="Custo"
          type="number"
          fullWidth
          variant="standard"
          value={formData.cost}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="stock_qty"
          label="Estoque"
          type="number"
          fullWidth
          variant="standard"
          value={formData.stock_qty}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="min_stock"
          label="Estoque Mínimo"
          type="number"
          fullWidth
          variant="standard"
          value={formData.min_stock}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="description"
          label="Descrição"
          type="text"
          fullWidth
          multiline
          rows={4}
          variant="standard"
          value={formData.description}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="tenant_id"
          label="Tenant ID"
          type="text"
          fullWidth
          variant="standard"
          value={formData.tenant_id}
          onChange={handleChange}
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
