import React, { useEffect, useState } from 'react';
import ClientForm from './ClientForm';
import { useNavigate, useParams } from 'react-router-dom';
import apiClient from '../../lib/apiClient';
import { Client } from './ClientListPage';

const ClientFormWrapper: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [client, setClient] = useState<Client | null>(null);

  useEffect(() => {
    if (id) {
      const fetchClient = async () => {
        try {
          const response = await apiClient<{ data: Client }>(`/admin/clients/${id}`);
          setClient(response.data);
        } catch (error) {
          console.error('Error fetching client:', error);
        }
      };
      fetchClient();
    }
  }, [id]);

  const handleClose = () => {
    navigate('/clients');
  };

  const handleSave = async (formData: Partial<Client>) => {
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
