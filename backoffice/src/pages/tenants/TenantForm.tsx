import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';

interface Tenant {
  id: string;
  name: string;
  document: string;
  phone: string;
  email: string;
}

interface TenantFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (tenant: Partial<Tenant>) => Promise<void>;
  tenant: Tenant | null;
}

const TenantForm: React.FC<TenantFormProps> = ({
  open,
  onClose,
  onSave,
  tenant,
}) => {
  const [formData, setFormData] = useState<Partial<Tenant>>({
    name: '',
    document: '',
    phone: '',
    email: '',
  });

  useEffect(() => {
    if (tenant) {
      setFormData(tenant);
    } else {
      setFormData({
        name: '',
        document: '',
        phone: '',
        email: '',
      });
    }
  }, [tenant, open]);

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
      <DialogTitle>{tenant ? 'Editar Tenant' : 'Adicionar Tenant'}</DialogTitle>
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
          name="document"
          label="Documento"
          type="text"
          fullWidth
          variant="standard"
          value={formData.document}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="phone"
          label="Telefone"
          type="text"
          fullWidth
          variant="standard"
          value={formData.phone}
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
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancelar</Button>
        <Button onClick={handleSave}>Salvar</Button>
      </DialogActions>
    </Dialog>
  );
};

export default TenantForm;
