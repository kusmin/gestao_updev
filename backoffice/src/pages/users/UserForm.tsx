import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';

interface User {
  id: string;
  name: string;
  email: string;
  phone: string;
  role: string;
  active: boolean;
  tenant_id: string;
}

interface UserFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (user: Partial<User>) => Promise<void>;
  user: User | null;
}

const UserForm: React.FC<UserFormProps> = ({
  open,
  onClose,
  onSave,
  user,
}) => {
  const [formData, setFormData] = useState<Partial<User>>({
    name: '',
    email: '',
    phone: '',
    role: 'professional',
    tenant_id: '',
  });

  useEffect(() => {
    if (user) {
      setFormData(user);
    } else {
      setFormData({
        name: '',
        email: '',
        phone: '',
        role: 'professional',
        tenant_id: '',
      });
    }
  }, [user, open]);

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
      <DialogTitle>{user ? 'Editar Usuário' : 'Adicionar Usuário'}</DialogTitle>
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
          type="text"
          fullWidth
          variant="standard"
          value={formData.phone}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="role"
          label="Role"
          type="text"
          fullWidth
          variant="standard"
          value={formData.role}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="password"
          label="Senha"
          type="password"
          fullWidth
          variant="standard"
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

export default UserForm;
