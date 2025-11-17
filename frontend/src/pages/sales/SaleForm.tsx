import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
  Typography,
  Box,
} from '@mui/material';
import { useAuth } from '../../contexts/useAuth';

// TODO: Substituir pelos tipos reais e pelas funções da API
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

const createSale = async (params: { tenantId: string; input: SaleRequest; accessToken: string }): Promise<Sale> => {
  console.log('Creating sale with params:', params);
  const total = params.input.items.reduce((acc, item) => acc + item.quantity * 1, 0); // Preço fixo de 1 para simplicidade
  return {
    id: new Date().toISOString(),
    clientName: params.input.clientName,
    date: new Date().toISOString(),
    total,
    items: params.input.items.map(item => ({ ...item, productName: `Produto ${item.productId}`, price: 1 })),
  };
};

const updateSale = async (params: { tenantId: string; saleId: string; input: SaleRequest; accessToken: string }): Promise<Sale> => {
  console.log('Updating sale with params:', params);
  const total = params.input.items.reduce((acc, item) => acc + item.quantity * 1, 0);
  return {
    id: params.saleId,
    clientName: params.input.clientName,
    date: new Date().toISOString(),
    total,
    items: params.input.items.map(item => ({ ...item, productName: `Produto ${item.productId}`, price: 1 })),
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
  // TODO: Implementar a lógica de adicionar/remover itens
  const [items, setItems] = useState<{ productId: string; quantity: number }[]>([]);
  const { tenantId, accessToken } = useAuth();

  useEffect(() => {
    if (sale) {
      setClientName(sale.clientName);
      setItems(sale.items.map(i => ({ productId: i.productId, quantity: i.quantity })));
    } else {
      setClientName('');
      setItems([]);
    }
  }, [sale]);

  const handleSave = async () => {
    if (!tenantId || !accessToken) {
      return;
    }
    const saleData: SaleRequest = { clientName, items };
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
    } catch (error) {
      console.error('Error saving sale:', error);
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
          {/* TODO: Implementar a UI para adicionar e listar itens */}
          <Typography>A funcionalidade de adicionar itens será implementada aqui.</Typography>
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
