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

import { Service } from '@/types/service';

interface ServiceFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (service: Partial<Service>) => Promise<void>;
  service: Service | null;
}

const ServiceForm: React.FC<ServiceFormProps> = ({ open, onClose, onSave, service }) => {
  const { control, handleSubmit, reset } = useForm<Service>({
    defaultValues: {
      name: '',
      category: '',
      description: '',
      duration_minutes: 30,
      price: 0,
      color: '#FFFFFF',
      tenant_id: '',
    },
  });

  useEffect(() => {
    if (open) {
      reset(service || {
        name: '',
        category: '',
        description: '',
        duration_minutes: 30,
        price: 0,
        color: '#FFFFFF',
        tenant_id: '',
      });
    }
  }, [service, open, reset]);

  const handleSave = (data: Service) => {
    onSave(data).then(() => {
      onClose();
    });
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{service ? 'Editar Serviço' : 'Adicionar Serviço'}</DialogTitle>
      <form onSubmit={handleSubmit(handleSave)}>
        <DialogContent>
          <Controller
            name="name"
            control={control}
            render={({ field }) => (
              <TextField {...field} autoFocus margin="dense" label="Nome" type="text" fullWidth variant="standard" />
            )}
          />
          <Controller
            name="category"
            control={control}
            render={({ field }) => (
              <TextField {...field} margin="dense" label="Categoria" type="text" fullWidth variant="standard" />
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
            name="duration_minutes"
            control={control}
            render={({ field }) => (
              <TextField
                {...field}
                margin="dense"
                label="Duração (minutos)"
                type="number"
                fullWidth
                variant="standard"
              />
            )}
          />
          <Controller
            name="price"
            control={control}
            render={({ field }) => (
              <TextField {...field} margin="dense" label="Preço" type="number" fullWidth variant="standard" />
            )}
          />
          <Controller
            name="color"
            control={control}
            render={({ field }) => (
              <TextField {...field} margin="dense" label="Cor" type="color" fullWidth variant="standard" />
            )}
          />
          <Controller
            name="tenant_id"
            control={control}
            render={({ field }) => (
              <TextField {...field} margin="dense" label="Tenant ID" type="text" fullWidth variant="standard" />
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

export default ServiceForm;
