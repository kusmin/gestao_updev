import React, { useState } from 'react';
import ProductForm from './ProductForm';
import { useNavigate, useParams } from 'react-router-dom';
import apiClient from '../../lib/apiClient';

const ProductFormWrapper: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [product, setProduct] = useState(null); // In a real app, fetch product by id

  const handleClose = () => {
    navigate('/products');
  };

  const handleSave = async (formData: any) => {
    // In a real app, handle save logic here
    console.log('Saving product:', formData);
    if (id) {
      await apiClient(`/admin/products/${id}`, { method: 'PUT', body: JSON.stringify(formData) });
    } else {
      await apiClient('/admin/products', { method: 'POST', body: JSON.stringify(formData) });
    }
    navigate('/products');
  };

  return (
    <ProductForm
      open={true} // Always open when rendered via route
      onClose={handleClose}
      onSave={handleSave}
      product={product}
    />
  );
};

export default ProductFormWrapper;
