import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';

interface Service {
  id: string;
  name: string;
  category: string;
  description: string;
  duration_minutes: number;
  price: number;
  color: string;
  tenant_id: string;
}

interface ServiceFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (service: Partial<Service>) => Promise<void>;
  service: Service | null;
}

const ServiceForm: React.FC<ServiceFormProps> = ({
  open,
  onClose,
  onSave,
  service,
}) => {
  const [formData, setFormData] = useState<Partial<Service>>({
    name: '',
    category: '',
    description: '',
    duration_minutes: 30,
    price: 0,
    color: '#FFFFFF',
    tenant_id: '',
  });

  useEffect(() => {
    if (service) {
      setFormData(service);
    } else {
      setFormData({
        name: '',
        category: '',
        description: '',
        duration_minutes: 30,
        price: 0,
        color: '#FFFFFF',
        tenant_id: '',
      });
    }
  }, [service, open]);

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
      <DialogTitle>{service ? 'Editar Serviço' : 'Adicionar Serviço'}</DialogTitle>
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
          name="category"
          label="Categoria"
          type="text"
          fullWidth
          variant="standard"
          value={formData.category}
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
          name="duration_minutes"
          label="Duração (minutos)"
          type="number"
          fullWidth
          variant="standard"
          value={formData.duration_minutes}
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
          name="color"
          label="Cor"
          type="color"
          fullWidth
          variant="standard"
          value={formData.color}
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

export default ServiceForm;
