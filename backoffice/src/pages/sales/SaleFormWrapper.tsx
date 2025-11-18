import React, { useState } from 'react';
import SaleForm from './SaleForm';
import { useNavigate, useParams } from 'react-router-dom';
import apiClient from '../../lib/apiClient';

const SaleFormWrapper: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [sale, setSale] = useState(null); // In a real app, fetch sale by id

  const handleClose = () => {
    navigate('/sales');
  };

  const handleSave = async (formData: any) => {
    // In a real app, handle save logic here
    console.log('Saving sale:', formData);
    if (id) {
      await apiClient(`/admin/sales/orders/${id}`, { method: 'PUT', body: JSON.stringify(formData) });
    } else {
      await apiClient('/admin/sales/orders', { method: 'POST', body: JSON.stringify(formData) });
    }
    navigate('/sales');
  };

  return (
    <SaleForm
      open={true} // Always open when rendered via route
      onClose={handleClose}
      onSave={handleSave}
      sale={sale}
    />
  );
};

export default SaleFormWrapper;
