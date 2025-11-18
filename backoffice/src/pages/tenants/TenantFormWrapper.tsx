import React, { useState } from 'react';
import TenantForm from './TenantForm';
import { useNavigate, useParams } from 'react-router-dom';
import apiClient from '../../lib/apiClient';

const TenantFormWrapper: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [tenant, setTenant] = useState(null); // In a real app, fetch tenant by id

  const handleClose = () => {
    navigate('/tenants');
  };

  const handleSave = async (formData: any) => {
    // In a real app, handle save logic here
    console.log('Saving tenant:', formData);
    if (id) {
      await apiClient(`/admin/tenants/${id}`, { method: 'PUT', body: JSON.stringify(formData) });
    } else {
      await apiClient('/admin/tenants', { method: 'POST', body: JSON.stringify(formData) });
    }
    navigate('/tenants');
  };

  return (
    <TenantForm
      open={true} // Always open when rendered via route
      onClose={handleClose}
      onSave={handleSave}
      tenant={tenant}
    />
  );
};

export default TenantFormWrapper;
