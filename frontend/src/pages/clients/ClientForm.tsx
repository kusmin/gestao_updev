import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';
import type { components } from '../../types/api';
import { createClient, updateClient } from '../../lib/apiClient';

type Client = components['schemas']['Client'];
type ClientRequest = components['schemas']['ClientRequest'];

interface ClientFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (client: any) => void;
  client: Client | null;
}

const ClientForm: React.FC<ClientFormProps> = ({
  open,
  onClose,
  onSave,
  client,
}) => {
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
    const clientData: ClientRequest = { name, email, phone };
    try {
      let savedClient;
      // TODO: Replace with actual tenant ID
      const tenantId = 'a4b2b2b2-b2b2-4b2b-b2b2-b2b2b2b2b2b2';
      if (client) {
        savedClient = await updateClient(tenantId, client.id, clientData);
      } else {
        savedClient = await createClient(tenantId, clientData);
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