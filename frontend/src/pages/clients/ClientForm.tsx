import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';
import { useMutation } from '@tanstack/react-query';
import { createClient, updateClient, type Client, type ClientRequest } from '../../lib/apiClient';
import { useAuth } from '../../contexts/useAuth';
import { useSnackbar } from '../../hooks/useSnackbar';

interface ClientFormProps {
  open: boolean;
  onClose: () => void;
  onSave: () => void;
  client: Client | null;
}

const ClientForm: React.FC<ClientFormProps> = ({ open, onClose, onSave, client }) => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [phone, setPhone] = useState('');
  const { tenantId, accessToken } = useAuth();
  const { showSnackbar } = useSnackbar();

  const createMutation = useMutation({
    mutationFn: (newClient: ClientRequest) => createClient({ tenantId: tenantId!, input: newClient, accessToken: accessToken! }),
    onSuccess: () => {
      onSave();
      onClose();
      showSnackbar('Cliente criado com sucesso!', 'success');
    },
    onError: () => {
      showSnackbar('Erro ao criar cliente.', 'error');
    },
  });

  const updateMutation = useMutation({
    mutationFn: (updatedClient: { clientId: string; input: ClientRequest }) => updateClient({ tenantId: tenantId!, ...updatedClient, accessToken: accessToken! }),
    onSuccess: () => {
      onSave();
      onClose();
      showSnackbar('Cliente atualizado com sucesso!', 'success');
    },
    onError: () => {
      showSnackbar('Erro ao atualizar cliente.', 'error');
    },
  });

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
  }, [client, open]);

  const handleSave = () => {
    if (!tenantId || !accessToken) {
      return;
    }
    const clientData: ClientRequest = { name, email, phone };
    if (client) {
      updateMutation.mutate({ clientId: client.id, input: clientData });
    } else {
      createMutation.mutate(clientData);
    }
  };

  const isLoading = createMutation.isPending || updateMutation.isPending;

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
        <Button onClick={onClose} disabled={isLoading}>Cancelar</Button>
        <Button onClick={handleSave} disabled={isLoading}>
          {isLoading ? 'Salvando...' : 'Salvar'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default ClientForm;
