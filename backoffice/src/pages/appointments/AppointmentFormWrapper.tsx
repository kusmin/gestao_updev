import React, { useEffect, useState } from 'react';
import AppointmentForm from './AppointmentForm';
import { useNavigate, useParams } from 'react-router-dom';
import apiClient from '../../lib/apiClient';
import { Appointment } from './AppointmentListPage';

const AppointmentFormWrapper: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [appointment, setAppointment] = useState<Appointment | null>(null);

  useEffect(() => {
    if (id) {
      const fetchAppointment = async () => {
        try {
          const response = await apiClient<{ data: Appointment }>(`/admin/bookings/${id}`);
          setAppointment(response.data);
        } catch (error) {
          console.error('Error fetching appointment:', error);
        }
      };
      fetchAppointment();
    }
  }, [id]);

  const handleClose = () => {
    navigate('/appointments');
  };

  const handleSave = async (formData: Partial<Appointment>) => {
    // In a real app, handle save logic here
    console.log('Saving appointment:', formData);
    if (id) {
      await apiClient(`/admin/bookings/${id}`, { method: 'PUT', body: JSON.stringify(formData) });
    } else {
      await apiClient('/admin/bookings', { method: 'POST', body: JSON.stringify(formData) });
    }
    navigate('/appointments');
  };

  return (
    <AppointmentForm
      open={true} // Always open when rendered via route
      onClose={handleClose}
      onSave={handleSave}
      appointment={appointment}
    />
  );
};

export default AppointmentFormWrapper;
