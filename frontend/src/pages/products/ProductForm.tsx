import React, { useEffect, useState } from 'react';
import {
  Button,
  DialogActions,
  TextField,
  Container,
  Typography,
  Box,
  Paper,
} from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';

const ProductForm: React.FC = () => {
  const [name, setName] = useState('');
  const [price, setPrice] = useState('');
  const [stock, setStock] = useState('');
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const isEditMode = Boolean(id);

  useEffect(() => {
    if (isEditMode) {
      // In a real app, you would fetch the product data here
      setName(`Produto ${id}`);
      setPrice('15.00');
      setStock('50');
    }
  }, [id, isEditMode]);

  const handleSave = () => {
    // In a real app, you would save the data to the backend
    alert(`Salvando produto: ${name}`);
    navigate('/products');
  };

  return (
    <Container maxWidth="sm">
      <Paper sx={{ p: 4, mt: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {isEditMode ? 'Editar Produto' : 'Adicionar Produto'}
        </Typography>
        <Box component="form" noValidate autoComplete="off">
          <TextField
            autoFocus
            margin="dense"
            id="name"
            label="Nome"
            type="text"
            fullWidth
            variant="outlined"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <TextField
            margin="dense"
            id="price"
            label="Preço"
            type="number"
            fullWidth
            variant="outlined"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
          />
          <TextField
            margin="dense"
            id="stock"
            label="Estoque"
            type="number"
            fullWidth
            variant="outlined"
            value={stock}
            onChange={(e) => setStock(e.target.value)}
          />
        </Box>
        <DialogActions sx={{ mt: 2, px: 0 }}>
          <Button onClick={() => navigate('/products')}>
            Cancelar
          </Button>
          <Button onClick={handleSave} variant="contained">
            Salvar
          </Button>
        </DialogActions>
      </Paper>
    </Container>
  );
};

export default ProductForm;