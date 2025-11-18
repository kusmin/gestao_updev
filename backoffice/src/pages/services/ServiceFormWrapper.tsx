import React, { useState } from 'react';
import ServiceForm from './ServiceForm';
import { useNavigate, useParams } from 'react-router-dom';
import apiClient from '../../lib/apiClient';

const ServiceFormWrapper: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [service, setService] = useState(null); // In a real app, fetch service by id

  const handleClose = () => {
    navigate('/services');
  };

  const handleSave = async (formData: any) => {
    // In a real app, handle save logic here
    console.log('Saving service:', formData);
    if (id) {
      await apiClient(`/admin/services/${id}`, { method: 'PUT', body: JSON.stringify(formData) });
    } else {
      await apiClient('/admin/services', { method: 'POST', body: JSON.stringify(formData) });
    }
    navigate('/services');
  };

  return (
    <ServiceForm
      open={true} // Always open when rendered via route
      onClose={handleClose}
      onSave={handleSave}
      service={service}
    />
  );
};

export default ServiceFormWrapper;
