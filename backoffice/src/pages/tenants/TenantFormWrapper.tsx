import React, { useEffect, useState } from 'react';
import TenantForm from './TenantForm';
import { useNavigate, useParams } from 'react-router-dom';
import apiClient from '../../lib/apiClient';
import { Tenant } from './TenantListPage';

const TenantFormWrapper: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [tenant, setTenant] = useState<Tenant | null>(null);

  useEffect(() => {
    if (id) {
      const fetchTenant = async () => {
        try {
          const response = await apiClient<{ data: Tenant }>(`/admin/tenants/${id}`);
          setTenant(response.data);
        } catch (error) {
          console.error('Error fetching tenant:', error);
        }
      };
      fetchTenant();
    }
  }, [id]);

  const handleClose = () => {
    navigate('/tenants');
  };

  const handleSave = async (formData: Partial<Tenant>) => {
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
