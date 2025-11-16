import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';
import { createClient, updateClient, type Client, type ClientRequest } from '../../lib/apiClient';
import { useAuth } from '../../contexts/useAuth';

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
  const { tenantId, accessToken } = useAuth();

  useEffect(() => {
    if (client) {
      setName(client.name);
      setEmail(client.email ?? '');
      setPhone(client.phone ?? '');
    } else {
      setName('');
      setEmail('');
      setPhone('');
    }
  }, [client]);

  const handleSave = async () => {
    if (!tenantId || !accessToken) {
      return;
    }
    const clientData: ClientRequest = { name, email, phone };
    try {
      let savedClient;
      if (client) {
        savedClient = await updateClient({
          tenantId,
          clientId: client.id,
          input: clientData,
          accessToken,
        });
      } else {
        savedClient = await createClient({ tenantId, input: clientData, accessToken });
      }
      onSave(savedClient);
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
