import React, { useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';

interface ClientFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (client: any) => void;
}

const ClientForm: React.FC<ClientFormProps> = ({ open, onClose, onSave }) => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');

  const handleSave = async () => {
    try {
      // TODO: Replace with actual API endpoint and add authentication
      const response = await fetch('http://localhost:8080/v1/clients', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          // TODO: Replace with actual tenant ID
          'X-Tenant-ID': 'a4b2b2b2-b2b2-4b2b-b2b2-b2b2b2b2b2b2',
        },
        body: JSON.stringify({ name, email, phone }),
      });
      const newClient = await response.json();
      onSave(newClient.data);
      onClose();
    } catch (error) {
      console.error('Error creating client:', error);
    }
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>Adicionar Cliente</DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          id="name"
          label="Nome"
          type="text"
          fullWidth
          variant="standard"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <TextField
          margin="dense"
          id="email"
          label="Email"
          type="email"
          fullWidth
          variant="standard"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <TextField
          margin="dense"
          id="phone"
          label="Telefone"
          type="tel"
          fullWidth
          variant="standard"
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
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