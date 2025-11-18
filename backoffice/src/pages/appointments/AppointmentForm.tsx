import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
  Select,
  MenuItem,
  InputLabel,
  FormControl,
  SelectChangeEvent,
} from '@mui/material';

interface Appointment {
  id: string;
  client_id: string;
  professional_id: string;
  service_id: string;
  start_at: string;
  end_at: string;
  status: string;
  tenant_id: string;
}

interface AppointmentFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (appointment: Partial<Appointment>) => Promise<void>;
  appointment: Appointment | null;
}

const AppointmentForm: React.FC<AppointmentFormProps> = ({ open, onClose, onSave, appointment }) => {
  const [formData, setFormData] = useState<Partial<Appointment>>({
    client_id: '',
    professional_id: '',
    service_id: '',
    start_at: '',
    end_at: '',
    status: '',
    tenant_id: '',
  });

  useEffect(() => {
    if (appointment) {
      setFormData(appointment);
    } else {
      setFormData({
        client_id: '',
        professional_id: '',
        service_id: '',
        start_at: '',
        end_at: '',
        status: '',
        tenant_id: '',
      });
    }
  }, [appointment, open]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement> | SelectChangeEvent<string>) => {
    const { name, value } = event.target;
    setFormData((prev) => ({ ...prev, [name as string]: value }));
  };

  const handleSave = async () => {
    // TODO: Add validation
    await onSave(formData);
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{appointment ? 'Editar Agendamento' : 'Adicionar Agendamento'}</DialogTitle>
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
          label="Start At"
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
          label="End At"
          type="datetime-local"
          fullWidth
          variant="standard"
          value={formData.end_at}
          onChange={handleChange}
          InputLabelProps={{
            shrink: true,
          }}
        />
        <FormControl fullWidth margin="dense" variant="standard">
          <InputLabel id="status-label">Status</InputLabel>
          <Select
            labelId="status-label"
            id="status"
            name="status"
            value={formData.status}
            onChange={handleChange}
            label="Status"
          >
            <MenuItem value="pending">Pending</MenuItem>
            <MenuItem value="confirmed">Confirmed</MenuItem>
            <MenuItem value="done">Done</MenuItem>
            <MenuItem value="canceled">Canceled</MenuItem>
          </Select>
        </FormControl>
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
