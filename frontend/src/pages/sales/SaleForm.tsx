import React, { useEffect, useState } from 'react';
import {
  Button,
  Container,
  Typography,
  Box,
  Paper,
  TextField,
  DialogActions,
} from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';

const SaleForm: React.FC = () => {
  const [client, setClient] = useState('');
  const [total, setTotal] = useState('');
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const isEditMode = Boolean(id);

  useEffect(() => {
    if (isEditMode) {
      // In a real app, you would fetch the sale data here
      setClient(`Cliente ${id}`);
      setTotal('150.00');
    }
  }, [id, isEditMode]);

  const handleSave = () => {
    // In a real app, you would save the data to the backend
    alert(`Salvando venda para: ${client}`);
    navigate('/sales');
  };

  return (
    <Container maxWidth="sm">
      <Paper sx={{ p: 4, mt: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          {isEditMode ? 'Editar Venda' : 'Adicionar Venda'}
        </Typography>
        <Box component="form" noValidate autoComplete="off">
          <TextField
            autoFocus
            margin="dense"
            id="client"
            label="Cliente"
            type="text"
            fullWidth
            variant="outlined"
            value={client}
            onChange={(e) => setClient(e.target.value)}
          />
          <TextField
            margin="dense"
            id="total"
            label="Total"
            type="number"
            fullWidth
            variant="outlined"
            value={total}
            onChange={(e) => setTotal(e.target.value)}
          />
        </Box>
        <DialogActions sx={{ mt: 2, px: 0 }}>
          <Button onClick={() => navigate('/sales')}>
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

export default SaleForm;