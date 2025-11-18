import React, { useState } from 'react';
import AppointmentForm from './AppointmentForm';
import { useNavigate, useParams } from 'react-router-dom';
import apiClient from '../../lib/apiClient';

const AppointmentFormWrapper: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [appointment, setAppointment] = useState(null); // In a real app, fetch appointment by id

  const handleClose = () => {
    navigate('/appointments');
  };

  const handleSave = async (formData: any) => {
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
