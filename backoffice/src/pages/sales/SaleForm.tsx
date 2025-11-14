import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';

interface SalesItem {
  ref_id: string;
  type: string;
  quantity: number;
  unit_price: number;
}

interface SalesOrder {
  id: string;
  client_id: string;
  booking_id?: string;
  discount?: number;
  notes?: string;
  items: SalesItem[];
  tenant_id: string;
}

interface SaleFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (sale: Partial<SalesOrder>) => Promise<void>;
  sale: SalesOrder | null;
}

const SaleForm: React.FC<SaleFormProps> = ({
  open,
  onClose,
  onSave,
  sale,
}) => {
  const [formData, setFormData] = useState<Partial<SalesOrder>>({
    client_id: '',
    items: [],
    tenant_id: '',
  });

  useEffect(() => {
    if (sale) {
      setFormData(sale);
    } else {
      setFormData({
        client_id: '',
        items: [],
        tenant_id: '',
      });
    }
  }, [sale, open]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = event.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    // TODO: Add validation and proper items handling
    await onSave(formData);
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{sale ? 'Editar Venda' : 'Adicionar Venda'}</DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          name="client_id"
          label="Client ID"
          type="text"
          fullWidth
          variant="standard"
          value={formData.client_id}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="booking_id"
          label="Booking ID"
          type="text"
          fullWidth
          variant="standard"
          value={formData.booking_id}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="discount"
          label="Desconto"
          type="number"
          fullWidth
          variant="standard"
          value={formData.discount}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="notes"
          label="Notas"
          type="text"
          fullWidth
          multiline
          rows={4}
          variant="standard"
          value={formData.notes}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="items"
          label="Itens (JSON)"
          type="text"
          fullWidth
          multiline
          rows={4}
          variant="standard"
          value={JSON.stringify(formData.items)}
          onChange={(e) =>
            setFormData((prev) => ({
              ...prev,
              items: JSON.parse(e.target.value),
            }))
          }
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

export default SaleForm;
