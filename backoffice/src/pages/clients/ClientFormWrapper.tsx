import React, { useState } from 'react';
import ClientForm from './ClientForm';
import { useNavigate, useParams } from 'react-router-dom';
import apiClient from '../../lib/apiClient';

const ClientFormWrapper: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [client, setClient] = useState(null); // In a real app, fetch client by id

  const handleClose = () => {
    navigate('/clients');
  };

  const handleSave = async (formData: any) => {
    // In a real app, handle save logic here
    console.log('Saving client:', formData);
    if (id) {
      await apiClient(`/admin/clients/${id}`, { method: 'PUT', body: JSON.stringify(formData) });
    } else {
      await apiClient('/admin/clients', { method: 'POST', body: JSON.stringify(formData) });
    }
    navigate('/clients');
  };

  return (
    <ClientForm
      open={true} // Always open when rendered via route
      onClose={handleClose}
      onSave={handleSave}
      client={client}
    />
  );
};

export default ClientFormWrapper;
