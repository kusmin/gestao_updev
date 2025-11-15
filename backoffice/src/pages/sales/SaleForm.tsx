import React, { useEffect } from 'react';
import { useForm, Controller } from 'react-hook-form';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';

import { SalesOrder } from '@/types/sales';

interface SaleFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (sale: Partial<SalesOrder>) => Promise<void>;
  sale: SalesOrder | null;
}

const SaleForm: React.FC<SaleFormProps> = ({ open, onClose, onSave, sale }) => {
  const { control, handleSubmit, reset } = useForm<SalesOrder>({
    defaultValues: {
      client_id: '',
      items: [],
      tenant_id: '',
    },
  });

  useEffect(() => {
    if (open) {
      reset(
        sale || {
          client_id: '',
          items: [],
          tenant_id: '',
        },
      );
    }
  }, [sale, open, reset]);

  const handleSave = (data: SalesOrder) => {
    onSave(data).then(() => {
      onClose();
    });
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{sale ? 'Editar Venda' : 'Adicionar Venda'}</DialogTitle>
      <form onSubmit={handleSubmit(handleSave)}>
        <DialogContent>
          <Controller
            name="client_id"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                autoFocus
                margin="dense"
                label="Client ID"
                type="text"
                fullWidth
                variant="standard"
              />
            )}
          />
          <Controller
            name="booking_id"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                margin="dense"
                label="Booking ID"
                type="text"
                fullWidth
                variant="standard"
              />
            )}
          />
          <Controller
            name="discount"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                margin="dense"
                label="Desconto"
                type="number"
                fullWidth
                variant="standard"
              />
            )}
          />
          <Controller
            name="notes"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                margin="dense"
                label="Notas"
                type="text"
                fullWidth
                multiline
                rows={4}
                variant="standard"
              />
            )}
          />
          <Controller
            name="items"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                margin="dense"
                label="Itens (JSON)"
                type="text"
                fullWidth
                multiline
                rows={4}
                variant="standard"
                value={JSON.stringify(field.value)}
                onChange={(e) => field.onChange(JSON.parse(e.target.value))}
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

export default SaleForm;
