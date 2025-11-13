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
}

interface ClientFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (client: Client) => void;
  client: Client | null;
}

const ClientForm: React.FC<ClientFormProps> = ({ open, onClose, onSave, client }) => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');

  useEffect(() => {
    if (client) {
      setName(client.name);
      setEmail(client.email);
      setPhone(client.phone);
    } else {
      setName('');
      setEmail('');
      setPhone('');
    }
  }, [client]);

  const handleSave = async () => {
    const clientData = { name, email, phone };
    try {
      const url = client
        ? `http://localhost:8080/v1/clients/${client.id}`
        : 'http://localhost:8080/v1/clients';
      const method = client ? 'PUT' : 'POST';

      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
          // TODO: Replace with actual tenant ID
          'X-Tenant-ID': 'a4b2b2b2-b2b2-4b2b-b2b2-b2b2b2b2b2b2',
        },
        body: JSON.stringify(clientData),
      });
      const savedClient = await response.json();
      onSave(savedClient.data);
      onClose();
    } catch (error) {
      console.error('Error saving client:', error);
    }
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{client ? 'Editar Cliente' : 'Adicionar Cliente'}</DialogTitle>
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
