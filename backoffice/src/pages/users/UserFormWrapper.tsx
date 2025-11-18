import React, { useState } from 'react';
import UserForm from './UserForm';
import { useNavigate, useParams } from 'react-router-dom';
import apiClient from '../../lib/apiClient';

const UserFormWrapper: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [user, setUser] = useState(null); // In a real app, fetch user by id

  const handleClose = () => {
    navigate('/users');
  };

  const handleSave = async (formData: any) => {
    // In a real app, handle save logic here
    console.log('Saving user:', formData);
    if (id) {
      await apiClient(`/admin/users/${id}`, { method: 'PUT', body: JSON.stringify(formData) });
    } else {
      await apiClient('/admin/users', { method: 'POST', body: JSON.stringify(formData) });
    }
    navigate('/users');
  };

  return (
    <UserForm
      open={true} // Always open when rendered via route
      onClose={handleClose}
      onSave={handleSave}
      user={user}
    />
  );
};

export default UserFormWrapper;
