import React, { useEffect, useState, useMemo } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
  Typography,
  Box,
  Autocomplete,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  IconButton,
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import { useAuth } from '../../contexts/useAuth';
import { useSnackbar } from '../../hooks/useSnackbar';

// TODO: Substituir pelos tipos reais e pelas funções da API
interface Product {
  id: string;
  name: string;
  price: number;
  stock: number;
}

interface Sale {
  id: string;
  clientName: string;
  total: number;
  date: string;
  items: SaleItem[];
}

interface SaleItem {
  productId: string;
  productName: string;
  quantity: number;
  price: number;
}

interface SaleRequest {
  clientName: string;
  items: { productId: string; quantity: number }[];
}

// Mock data e funções de API (serão substituídas pela integração real)
const mockProducts: Product[] = [
  { id: '1', name: 'Pomada Modeladora', price: 25.0, stock: 50 },
  { id: '2', name: 'Óleo para Barba', price: 35.0, stock: 30 },
  { id: '3', name: 'Shampoo para Cabelo', price: 20.0, stock: 100 },
];

const fetchProducts = async (): Promise<Product[]> => {
  return Promise.resolve(mockProducts);
};


const createSale = async (params: { tenantId: string; input: SaleRequest; accessToken: string }): Promise<Sale> => {
  console.log('Creating sale with params:', params);
  const total = params.input.items.reduce((acc, item) => {
    const product = mockProducts.find(p => p.id === item.productId);
    return acc + item.quantity * (product?.price || 0);
  }, 0);
  return {
    id: new Date().toISOString(),
    clientName: params.input.clientName,
    date: new Date().toISOString(),
    total,
    items: params.input.items.map(item => {
      const product = mockProducts.find(p => p.id === item.productId);
      return { ...item, productName: product?.name || '', price: product?.price || 0 };
    }),
  };
};

const updateSale = async (params: { tenantId: string; saleId: string; input: SaleRequest; accessToken: string }): Promise<Sale> => {
  console.log('Updating sale with params:', params);
  const total = params.input.items.reduce((acc, item) => {
    const product = mockProducts.find(p => p.id === item.productId);
    return acc + item.quantity * (product?.price || 0);
  }, 0);
  return {
    id: params.saleId,
    clientName: params.input.clientName,
    date: new Date().toISOString(),
    total,
    items: params.input.items.map(item => {
      const product = mockProducts.find(p => p.id === item.productId);
      return { ...item, productName: product?.name || '', price: product?.price || 0 };
    }),
  };
};

interface SaleFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (sale: Sale) => void;
  sale: Sale | null;
}

const SaleForm: React.FC<SaleFormProps> = ({ open, onClose, onSave, sale }) => {
  const [clientName, setClientName] = useState('');
  const [items, setItems] = useState<SaleItem[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [quantity, setQuantity] = useState(1);
  const { tenantId, accessToken } = useAuth();
  const { showSnackbar } = useSnackbar();

  useEffect(() => {
    fetchProducts().then(setProducts);
  }, []);

  useEffect(() => {
    if (sale) {
      setClientName(sale.clientName);
      setItems(sale.items);
    } else {
      setClientName('');
      setItems([]);
    }
  }, [sale]);

  const total = useMemo(() => {
    return items.reduce((acc, item) => acc + item.price * item.quantity, 0);
  }, [items]);

  const handleAddItem = () => {
    if (selectedProduct && quantity > 0) {
      const existingItem = items.find(item => item.productId === selectedProduct.id);
      if (existingItem) {
        setItems(items.map(item => item.productId === selectedProduct.id ? { ...item, quantity: item.quantity + quantity } : item));
      } else {
        setItems([...items, {
          productId: selectedProduct.id,
          productName: selectedProduct.name,
          quantity,
          price: selectedProduct.price,
        }]);
      }
      setSelectedProduct(null);
      setQuantity(1);
    }
  };

  const handleRemoveItem = (productId: string) => {
    setItems(items.filter(item => item.productId !== productId));
  };

  const handleSave = async () => {
    if (!tenantId || !accessToken) {
      return;
    }
    const saleData: SaleRequest = { clientName, items: items.map(i => ({ productId: i.productId, quantity: i.quantity })) };
    try {
      let savedSale;
      if (sale) {
        savedSale = await updateSale({
          tenantId,
          saleId: sale.id,
          input: saleData,
          accessToken,
        });
      } else {
        savedSale = await createSale({ tenantId, input: saleData, accessToken });
      }
      onSave(savedSale);
      onClose();
      showSnackbar(`Venda ${sale ? 'atualizada' : 'registrada'} com sucesso!`, 'success');
    } catch (error) {
      console.error('Error saving sale:', error);
      showSnackbar('Erro ao salvar venda.', 'error');
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>{sale ? 'Detalhes da Venda' : 'Registrar Venda'}</DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          id="clientName"
          label="Nome do Cliente"
          type="text"
          fullWidth
          variant="standard"
          value={clientName}
          onChange={(e) => setClientName(e.target.value)}
        />
        <Box sx={{ my: 2 }}>
          <Typography variant="h6">Itens da Venda</Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, my: 2 }}>
            <Autocomplete
              options={products}
              getOptionLabel={(option) => option.name}
              value={selectedProduct}
              onChange={(_, newValue) => setSelectedProduct(newValue)}
              renderInput={(params) => <TextField {...params} label="Produto" sx={{ flex: 1 }} />}
            />
            <TextField
              label="Qtd."
              type="number"
              value={quantity}
              onChange={(e) => setQuantity(Number(e.target.value))}
              sx={{ width: '100px' }}
            />
            <Button variant="contained" onClick={handleAddItem} disabled={!selectedProduct}>Adicionar</Button>
          </Box>
          <TableContainer>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Produto</TableCell>
                  <TableCell align="right">Qtd.</TableCell>
                  <TableCell align="right">Preço Unit.</TableCell>
                  <TableCell align="right">Subtotal</TableCell>
                  <TableCell align="right">Ações</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {items.map((item) => (
                  <TableRow key={item.productId}>
                    <TableCell>{item.productName}</TableCell>
                    <TableCell align="right">{item.quantity}</TableCell>
                    <TableCell align="right">R$ {item.price.toFixed(2)}</TableCell>
                    <TableCell align="right">R$ {(item.price * item.quantity).toFixed(2)}</TableCell>
                    <TableCell align="right">
                      <IconButton onClick={() => handleRemoveItem(item.productId)} size="small">
                        <DeleteIcon />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
          <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
            <Typography variant="h6">Total: R$ {total.toFixed(2)}</Typography>
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancelar</Button>
        <Button onClick={handleSave} disabled={items.length === 0}>Salvar</Button>
      </DialogActions>
    </Dialog>
  );
};

export default SaleForm;
