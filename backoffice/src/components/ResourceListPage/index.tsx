import React, { useEffect, useState, ReactNode, useCallback } from 'react';
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
} from '@mui/material';
import apiClient from '@/lib/apiClient';

interface Resource {
  id: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  [key: string]: any;
}

interface ResourceListPageProps<T extends Resource> {
  title: string;
  endpoint: string;
  columns: {
    header: string;
    accessor: keyof T;
  }[];
  renderForm: (props: {
    open: boolean;
    onClose: () => void;
    onSave: (item: Partial<T>) => Promise<void>;
    item: T | null;
  }) => ReactNode;
}

const ResourceListPage = <T extends Resource>({
  title,
  endpoint,
  columns,
  renderForm,
}: ResourceListPageProps<T>) => {
  const [items, setItems] = useState<T[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<T | null>(null);

  const fetchItems = useCallback(async () => {
    try {
      const response = await apiClient<{ data: T[] }>(endpoint);
      setItems(response.data);
    } catch (error) {
      console.error(`Error fetching ${title}:`, error);
    }
  }, [endpoint, title]);

  useEffect(() => {
    fetchItems();
  }, [fetchItems]);

  const handleOpenForm = (item: T | null = null) => {
    setEditingItem(item);
    setIsFormOpen(true);
  };

  const handleCloseForm = () => {
    setEditingItem(null);
    setIsFormOpen(false);
    fetchItems();
  };

  const handleSaveItem = async (item: Partial<T>) => {
    try {
      if (editingItem) {
        await apiClient(`${endpoint}/${editingItem.id}`, {
          method: 'PUT',
          body: JSON.stringify(item),
        });
      } else {
        await apiClient(endpoint, {
          method: 'POST',
          body: JSON.stringify(item),
        });
      }
    } catch (error) {
      console.error(`Error saving ${title}:`, error);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await apiClient(`${endpoint}/${id}`, { method: 'DELETE' });
      fetchItems();
    } catch (error) {
      console.error(`Error deleting ${title}:`, error);
    }
  };

  return (
    <Container>
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {title}
        </Typography>
        <Button variant="contained" color="primary" onClick={() => handleOpenForm()}>
          Adicionar {title}
        </Button>
      </Box>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              {columns.map((col) => (
                <TableCell key={col.accessor as string}>{col.header}</TableCell>
              ))}
              <TableCell>Ações</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {items.map((item) => (
              <TableRow key={item.id}>
                {columns.map((col) => (
                  <TableCell key={col.accessor as string}>{item[col.accessor]}</TableCell>
                ))}
                <TableCell>
                  <Button onClick={() => handleOpenForm(item)}>Editar</Button>
                  <Button color="error" onClick={() => handleDelete(item.id)}>
                    Excluir
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      {renderForm({
        open: isFormOpen,
        onClose: handleCloseForm,
        onSave: handleSaveItem,
        item: editingItem,
      })}
    </Container>
  );
};

export default ResourceListPage;
