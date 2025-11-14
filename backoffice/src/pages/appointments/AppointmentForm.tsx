import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';

interface Booking {
  id: string;
  client_id: string;
  professional_id: string;
  service_id: string;
  status: string;
  start_at: string;
  end_at: string;
  tenant_id: string;
}

interface AppointmentFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (booking: Partial<Booking>) => Promise<void>;
  booking: Booking | null;
}

const AppointmentForm: React.FC<AppointmentFormProps> = ({ open, onClose, onSave, booking }) => {
  const [formData, setFormData] = useState<Partial<Booking>>({
    client_id: '',
    professional_id: '',
    service_id: '',
    status: 'pending',
    start_at: '',
    end_at: '',
    tenant_id: '',
  });

  useEffect(() => {
    if (booking) {
      setFormData(booking);
    } else {
      setFormData({
        client_id: '',
        professional_id: '',
        service_id: '',
        status: 'pending',
        start_at: '',
        end_at: '',
        tenant_id: '',
      });
    }
  }, [booking, open]);

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
      <DialogTitle>{booking ? 'Editar Agendamento' : 'Adicionar Agendamento'}</DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          name="client_id"
          label="Client ID"
          type="text"
          fullWidth
          variant="standard"
          value={formData.client_id}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="professional_id"
          label="Professional ID"
          type="text"
          fullWidth
          variant="standard"
          value={formData.professional_id}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="service_id"
          label="Service ID"
          type="text"
          fullWidth
          variant="standard"
          value={formData.service_id}
          onChange={handleChange}
        />
        <TextField
          margin="dense"
          name="start_at"
          label="Início"
          type="datetime-local"
          fullWidth
          variant="standard"
          value={formData.start_at}
          onChange={handleChange}
          InputLabelProps={{
            shrink: true,
          }}
        />
        <TextField
          margin="dense"
          name="end_at"
          label="Fim"
          type="datetime-local"
          fullWidth
          variant="standard"
          value={formData.end_at}
          onChange={handleChange}
          InputLabelProps={{
            shrink: true,
          }}
        />
        <TextField
          margin="dense"
          name="status"
          label="Status"
          type="text"
          fullWidth
          variant="standard"
          value={formData.status}
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

export default AppointmentForm;
