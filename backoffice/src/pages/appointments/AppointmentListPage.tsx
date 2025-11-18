import React, { useEffect, useState } from 'react';
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
import apiClient from '../../lib/apiClient';

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

const AppointmentListPage: React.FC = () => {
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingAppointment, setEditingAppointment] = useState<Appointment | null>(null);

  const fetchAppointments = async () => {
    try {
      const response = await apiClient<{ data: Appointment[] }>('/admin/bookings');
      setAppointments(response.data);
    } catch (error) {
      console.error('Error fetching appointments:', error);
    }
  };

  useEffect(() => {
    fetchAppointments();
  }, []);

  const handleOpenForm = (appointment: Appointment | null = null) => {
    setEditingAppointment(appointment);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingAppointment(null);
    setIsFormOpen(false);
    fetchAppointments();
  };

  const handleSaveAppointment = async (appointment: Partial<Appointment>) => {
    try {
      if (editingAppointment) {
        await apiClient(`/admin/bookings/${editingAppointment.id}`, {
          method: 'PUT',
          body: JSON.stringify(appointment),
        });
      } else {
        await apiClient('/admin/bookings', {
          method: 'POST',
          body: JSON.stringify(appointment),
        });
      }
    } catch (error) {
      console.error('Error saving appointment:', error);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiClient(`/admin/bookings/${id}`, { method: 'DELETE' });
      fetchAppointments();
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
              <TableCell>Cliente ID</TableCell>
              <TableCell>Profissional ID</TableCell>
              <TableCell>Serviço ID</TableCell>
              <TableCell>Início</TableCell>
              <TableCell>Fim</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Tenant ID</TableCell>
              <TableCell>Ações</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {appointments.map((appointment) => (
              <TableRow key={appointment.id}>
                <TableCell>{appointment.client_id}</TableCell>
                <TableCell>{appointment.professional_id}</TableCell>
                <TableCell>{appointment.service_id}</TableCell>
                <TableCell>{new Date(appointment.start_at).toLocaleString()}</TableCell>
                <TableCell>{new Date(appointment.end_at).toLocaleString()}</TableCell>
                <TableCell>{appointment.status}</TableCell>
                <TableCell>{appointment.tenant_id}</TableCell>
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