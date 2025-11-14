import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';

interface Client {
  id: string;
  name: string;
  email: string;
  phone: string;
  tenant_id: string;
}

interface ClientFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (client: Partial<Client>) => Promise<void>;
  client: Client | null;
}

const ClientForm: React.FC<ClientFormProps> = ({ open, onClose, onSave, client }) => {
  const [formData, setFormData] = useState<Partial<Client>>({
    name: '',
    email: '',
    phone: '',
    tenant_id: '',
  });

  useEffect(() => {
    if (client) {
      setFormData(client);
    } else {
      setFormData({
        name: '',
        email: '',
        phone: '',
        tenant_id: '',
      });
    }
  }, [client, open]);

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
      <DialogTitle>{client ? 'Editar Cliente' : 'Adicionar Cliente'}</DialogTitle>
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
          name="email"
          label="Email"
          type="email"
          fullWidth
          variant="standard"
          value={formData.email}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="phone"
          label="Telefone"
          type="tel"
          fullWidth
          variant="standard"
          value={formData.phone}
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

export default ClientForm;
