import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
  Select,
  MenuItem,
  InputLabel,
  FormControl,
} from '@mui/material';

interface Sale {
  id: string;
  client_id: string;
  total: number;
  status: string;
  tenant_id: string;
}

interface SaleFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (sale: Partial<Sale>) => Promise<void>;
  sale: Sale | null;
}

const SaleForm: React.FC<SaleFormProps> = ({ open, onClose, onSave, sale }) => {
  const [formData, setFormData] = useState<Partial<Sale>>({
    client_id: '',
    total: 0,
    status: '',
    tenant_id: '',
  });

  useEffect(() => {
    if (sale) {
      setFormData(sale);
    } else {
      setFormData({
        client_id: '',
        total: 0,
        status: '',
        tenant_id: '',
      });
    }
  }, [sale, open]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement | { name?: string; value: unknown }>) => {
    const { name, value } = event.target;
    setFormData((prev) => ({ ...prev, [name as string]: value }));
  };

  const handleSave = async () => {
    // TODO: Add validation
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
          name="total"
          label="Total"
          type="number"
          fullWidth
          variant="standard"
          value={formData.total}
          onChange={handleChange}
        />
        <FormControl fullWidth margin="dense" variant="standard">
          <InputLabel id="status-label">Status</InputLabel>
          <Select
            labelId="status-label"
            id="status"
            name="status"
            value={formData.status}
            onChange={handleChange}
            label="Status"
          >
            <MenuItem value="pending">Pending</MenuItem>
            <MenuItem value="completed">Completed</MenuItem>
            <MenuItem value="canceled">Canceled</MenuItem>
          </Select>
        </FormControl>
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