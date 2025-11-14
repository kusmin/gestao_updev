import React, { useEffect, useState } from 'react';
import {
  Box,
  Button,
  Container,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  Switch,
} from '@mui/material';
import UserForm from './UserForm';
import apiClient from '../../lib/apiClient';

interface User {
  id: string;
  name: string;
  email: string;
  phone: string;
  role: string;
  active: boolean;
  tenant_id: string;
}

const UserListPage: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  const fetchUsers = async () => {
    try {
      const response = await apiClient<{ data: User[] }>('/admin/users');
      setUsers(response.data);
    } catch (error) {
      console.error('Error fetching users:', error);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleOpenForm = (user: User | null = null) => {
    setEditingUser(user);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingUser(null);
    setIsFormOpen(false);
    fetchUsers(); // Refetch users after closing form
  };

  const handleSaveUser = async (user: Partial<User>) => {
    try {
      if (editingUser) {
        await apiClient(`/admin/users/${editingUser.id}`, {
          method: 'PUT',
          body: JSON.stringify(user),
        });
      } else {
        await apiClient('/admin/users', {
          method: 'POST',
          body: JSON.stringify(user),
        });
      }
    } catch (error) {
      console.error('Error saving user:', error);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiClient(`/admin/users/${id}`, { method: 'DELETE' });
      fetchUsers(); // Refetch users after deleting
    } catch (error) {
      console.error('Error deleting user:', error);
    }
  };

  const handleToggleActive = async (user: User) => {
    try {
      await apiClient(`/admin/users/${user.id}`, {
        method: 'PUT',
        body: JSON.stringify({ active: !user.active }),
      });
      fetchUsers();
    } catch (error) {
      console.error('Error updating user:', error);
    }
  };

  return (
    <Container>
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Users
        </Typography>
        <Button
          variant="contained"
          color="primary"
          onClick={() => handleOpenForm()}
        >
          Adicionar Usuário
        </Button>
      </Box>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Nome</TableCell>
              <TableCell>Email</TableCell>
              <TableCell>Telefone</TableCell>
              <TableCell>Role</TableCell>
              <TableCell>Ativo</TableCell>
              <TableCell>Ações</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {users.map((user) => (
              <TableRow key={user.id}>
                <TableCell>{user.name}</TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell>{user.phone}</TableCell>
                <TableCell>{user.role}</TableCell>
                <TableCell>
                  <Switch
                    checked={user.active}
                    onChange={() => handleToggleActive(user)}
                  />
                </TableCell>
                <TableCell>
                  <Button onClick={() => handleOpenForm(user)}>Editar</Button>
                  <Button color="error" onClick={() => handleDelete(user.id)}>
                    Excluir
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <UserForm
        open={isFormOpen}
        onClose={handleCloseForm}
        onSave={handleSaveUser}
        user={editingUser}
      />
    </Container>
  );
};

export default UserListPage;