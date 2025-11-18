import React, { useEffect, useState } from 'react';
import ProductForm from './ProductForm';
import { useNavigate, useParams } from 'react-router-dom';
import apiClient from '../../lib/apiClient';
import { Product } from './ProductListPage';

const ProductFormWrapper: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [product, setProduct] = useState<Product | null>(null);

  useEffect(() => {
    if (id) {
      const fetchProduct = async () => {
        try {
          const response = await apiClient<{ data: Product }>(`/admin/products/${id}`);
          setProduct(response.data);
        } catch (error) {
          console.error('Error fetching product:', error);
        }
      };
      fetchProduct();
    }
  }, [id]);

  const handleClose = () => {
    navigate('/products');
  };

  const handleSave = async (formData: Partial<Product>) => {
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
