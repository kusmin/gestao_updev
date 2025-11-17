import React, { useEffect, useState } from 'react';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from '@mui/material';
import { useAuth } from '../../contexts/useAuth';
import { useSnackbar } from '../../hooks/useSnackbar';

// TODO: Substituir pelo tipo Appointment real e pelas funções da API
interface Appointment {
  id: string;
  clientName: string;
  date: string;
  service: string;
}

interface AppointmentRequest {
  clientName: string;
  date: string;
  service: string;
}

const createAppointment = async (params: { tenantId: string; input: AppointmentRequest; accessToken: string }): Promise<Appointment> => {
  console.log('Creating appointment with params:', params);
  return { ...params.input, id: new Date().toISOString() };
};

const updateAppointment = async (params: { tenantId: string; appointmentId: string; input: AppointmentRequest; accessToken: string }): Promise<Appointment> => {
  console.log('Updating appointment with params:', params);
  return { ...params.input, id: params.appointmentId };
};


interface AppointmentFormProps {
  open: boolean;
  onClose: () => void;
  onSave: (appointment: Appointment) => void;
  appointment: Appointment | null;
}

const AppointmentForm: React.FC<AppointmentFormProps> = ({ open, onClose, onSave, appointment }) => {
  const [clientName, setClientName] = useState('');
  const [date, setDate] = useState('');
  const [service, setService] = useState('');
  const { tenantId, accessToken } = useAuth();
  const { showSnackbar } = useSnackbar();

  useEffect(() => {
    if (appointment) {
      setClientName(appointment.clientName);
      setDate(appointment.date);
      setService(appointment.service);
    } else {
      setClientName('');
      setDate('');
      setService('');
    }
  }, [appointment]);

  const handleSave = async () => {
    if (!tenantId || !accessToken) {
      return;
    }
    const appointmentData: AppointmentRequest = { clientName, date, service };
    try {
      let savedAppointment;
      if (appointment) {
        savedAppointment = await updateAppointment({
          tenantId,
          appointmentId: appointment.id,
          input: appointmentData,
          accessToken,
        });
      } else {
        savedAppointment = await createAppointment({ tenantId, input: appointmentData, accessToken });
      }
      onSave(savedAppointment);
      onClose();
      showSnackbar(`Agendamento ${appointment ? 'atualizado' : 'criado'} com sucesso!`, 'success');
    } catch (error) {
      console.error('Error saving appointment:', error);
      showSnackbar('Erro ao salvar agendamento.', 'error');
    }
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{appointment ? 'Editar Agendamento' : 'Adicionar Agendamento'}</DialogTitle>
      <DialogContent>
        <TextField
          autoFocus
          margin="dense"
          id="clientName"
          label="Nome do Cliente"
          type="text"
          fullWidth
          variant="standard"
          value={clientName}
          onChange={(e) => setClientName(e.target.value)}
        />
        <TextField
          margin="dense"
          id="date"
          label="Data"
          type="datetime-local"
          fullWidth
          variant="standard"
          value={date}
          onChange={(e) => setDate(e.target.value)}
          InputLabelProps={{
            shrink: true,
          }}
        />
        <TextField
          margin="dense"
          id="service"
          label="Serviço"
          type="text"
          fullWidth
          variant="standard"
          value={service}
          onChange={(e) => setService(e.target.value)}
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
