import React, { useEffect, useState } from 'react';
import UserForm from './UserForm';
import { useNavigate, useParams } from 'react-router-dom';
import apiClient from '../../lib/apiClient';
import { User } from './UserListPage';

const UserFormWrapper: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    if (id) {
      const fetchUser = async () => {
        try {
          const response = await apiClient<{ data: User }>(`/admin/users/${id}`);
          setUser(response.data);
        } catch (error) {
          console.error('Error fetching user:', error);
        }
      };
      fetchUser();
    }
  }, [id]);

  const handleClose = () => {
    navigate('/users');
  };

  const handleSave = async (formData: Partial<User>) => {
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
