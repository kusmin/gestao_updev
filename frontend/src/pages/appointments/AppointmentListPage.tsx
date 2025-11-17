import React, { useCallback, useEffect, useState } from 'react';
import {
  Box,
  Button,
  Container,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material';
import AppointmentForm from './AppointmentForm';
import { useAuth } from '../../contexts/useAuth';

// TODO: Substituir pelo tipo Appointment real e pelas funções da API
interface Appointment {
  id: string;
  clientName: string;
  date: string;
  service: string;
}

const fetchAppointments = async (params: { tenantId: string; accessToken: string }): Promise<Appointment[]> => {
  console.log('Fetching appointments with params:', params);
  // Placeholder data
  return [
    { id: '1', clientName: 'João da Silva', date: '2025-11-18T10:00:00Z', service: 'Corte de Cabelo' },
    { id: '2', clientName: 'Maria Oliveira', date: '2025-11-18T11:00:00Z', service: 'Manicure' },
  ];
};

const deleteAppointment = async (params: { tenantId: string; appointmentId: string; accessToken: string }): Promise<void> => {
  console.log('Deleting appointment with params:', params);
  return Promise.resolve();
};


const AppointmentListPage: React.FC = () => {
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingAppointment, setEditingAppointment] = useState<Appointment | null>(null);
  const { tenantId, accessToken } = useAuth();

  const getAppointments = useCallback(async () => {
    if (!tenantId || !accessToken) {
      setAppointments([]);
      return;
    }
    try {
      const data = await fetchAppointments({ tenantId, accessToken });
      if (data && Array.isArray(data)) {
        setAppointments(data);
      }
    } catch (error) {
      console.error('Error fetching appointments:', error);
    }
  }, [tenantId, accessToken, setAppointments]);

  useEffect(() => {
    if (tenantId && accessToken) {
      getAppointments();
    }
  }, [tenantId, accessToken, getAppointments]);

  const handleOpenForm = (appointment: Appointment | null = null) => {
    setEditingAppointment(appointment);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingAppointment(null);
    setIsFormOpen(false);
  };

  const handleSaveAppointment = (appointment: Appointment) => {
    if (editingAppointment) {
      setAppointments(appointments.map((a) => (a.id === appointment.id ? appointment : a)));
    } else {
      setAppointments([...appointments, appointment]);
    }
    handleCloseForm();
  };

  const handleDelete = async (id: string) => {
    if (!tenantId || !accessToken) {
      return;
    }
    try {
      await deleteAppointment({ tenantId, appointmentId: id, accessToken });
      setAppointments((prev) => prev.filter((appointment) => appointment.id !== id));
    } catch (error) {
      console.error('Error deleting appointment:', error);
    }
  };

  return (
    <Container>
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Agendamentos
        </Typography>
        <Button variant="contained" color="primary" onClick={() => handleOpenForm()}>
          Adicionar Agendamento
        </Button>
      </Box>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Cliente</TableCell>
              <TableCell>Data</TableCell>
              <TableCell>Serviço</TableCell>
              <TableCell>Ações</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {appointments.map((appointment) => (
              <TableRow key={appointment.id}>
                <TableCell>{appointment.clientName}</TableCell>
                <TableCell>{new Date(appointment.date).toLocaleString()}</TableCell>
                <TableCell>{appointment.service}</TableCell>
                <TableCell>
                  <Button onClick={() => handleOpenForm(appointment)}>Editar</Button>
                  <Button color="error" onClick={() => handleDelete(appointment.id)}>
                    Excluir
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <AppointmentForm
        open={isFormOpen}
        onClose={handleCloseForm}
        onSave={handleSaveAppointment}
        appointment={editingAppointment}
      />
    </Container>
  );
};

export default AppointmentListPage;
