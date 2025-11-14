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

const AppointmentListPage: React.FC = () => {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingBooking, setEditingBooking] = useState<Booking | null>(null);

  const fetchBookings = async () => {
    try {
      const response = await apiClient<{ data: Booking[] }>('/admin/bookings');
      setBookings(response.data);
    } catch (error) {
      console.error('Error fetching bookings:', error);
    }
  };

  useEffect(() => {
    fetchBookings();
  }, []);

  const handleOpenForm = (booking: Booking | null = null) => {
    setEditingBooking(booking);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingBooking(null);
    setIsFormOpen(false);
    fetchBookings(); // Refetch bookings after closing form
  };

  const handleSaveBooking = async (booking: Partial<Booking>) => {
    try {
      if (editingBooking) {
        await apiClient(`/admin/bookings/${editingBooking.id}`, {
          method: 'PUT',
          body: JSON.stringify(booking),
        });
      } else {
        await apiClient('/admin/bookings', {
          method: 'POST',
          body: JSON.stringify(booking),
        });
      }
    } catch (error) {
      console.error('Error saving booking:', error);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiClient(`/admin/bookings/${id}`, { method: 'DELETE' });
      fetchBookings(); // Refetch bookings after deleting
    } catch (error) {
      console.error('Error deleting booking:', error);
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
              <TableCell>Profissional</TableCell>
              <TableCell>Serviço</TableCell>
              <TableCell>Início</TableCell>
              <TableCell>Fim</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Ações</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {bookings.map((booking) => (
              <TableRow key={booking.id}>
                <TableCell>{booking.client_id}</TableCell>
                <TableCell>{booking.professional_id}</TableCell>
                <TableCell>{booking.service_id}</TableCell>
                <TableCell>{new Date(booking.start_at).toLocaleString()}</TableCell>
                <TableCell>{new Date(booking.end_at).toLocaleString()}</TableCell>
                <TableCell>{booking.status}</TableCell>
                <TableCell>
                  <Button onClick={() => handleOpenForm(booking)}>Editar</Button>
                  <Button color="error" onClick={() => handleDelete(booking.id)}>
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
        onSave={handleSaveBooking}
        booking={editingBooking}
      />
    </Container>
  );
};

export default AppointmentListPage;
