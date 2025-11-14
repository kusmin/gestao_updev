import React, { useEffect, useState } from 'react';
import { Box, Typography, Grid, Card, CardContent } from '@mui/material';
import apiClient from '../lib/apiClient';

interface OverallMetrics {
  total_tenants: number;
  total_users: number;
  total_clients: number;
  total_products: number;
  total_services: number;
  total_bookings: number;
  total_revenue: number;
}

const DashboardPage: React.FC = () => {
  const [metrics, setMetrics] = useState<OverallMetrics | null>(null);

  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        const response = await apiClient<{ data: OverallMetrics }>('/admin/dashboard/metrics');
        setMetrics(response.data);
      } catch (error) {
        console.error('Error fetching overall metrics:', error);
      }
    };

    fetchMetrics();
  }, []);

  return (
    <Box sx={{ my: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Dashboard
      </Typography>

      {metrics ? (
        <Grid container spacing={3}>
          <Grid item xs={12} sm={6} md={4} component="div">
            <Card>
              <CardContent>
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  Total de Tenants
                </Typography>
                <Typography variant="h4">{metrics.total_tenants}</Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={4} component="div">
            <Card>
              <CardContent>
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  Total de Usuários
                </Typography>
                <Typography variant="h4">{metrics.total_users}</Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={4} component="div">
            <Card>
              <CardContent>
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  Total de Clientes
                </Typography>
                <Typography variant="h4">{metrics.total_clients}</Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={4} component="div">
            <Card>
              <CardContent>
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  Total de Produtos
                </Typography>
                <Typography variant="h4">{metrics.total_products}</Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={4} component="div">
            <Card>
              <CardContent>
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  Total de Serviços
                </Typography>
                <Typography variant="h4">{metrics.total_services}</Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={4} component="div">
            <Card>
              <CardContent>
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  Total de Agendamentos
                </Typography>
                <Typography variant="h4">{metrics.total_bookings}</Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={4} component="div">
            <Card>
              <CardContent>
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  Receita Total
                </Typography>
                <Typography variant="h4">
                  {metrics.total_revenue.toLocaleString('pt-BR', {
                    style: 'currency',
                    currency: 'BRL',
                  })}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      ) : (
        <Typography>Carregando métricas...</Typography>
      )}
    </Box>
  );
};

export default DashboardPage;
