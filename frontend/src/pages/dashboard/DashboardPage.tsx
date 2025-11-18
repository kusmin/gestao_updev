import { Card, CardContent, Grid, Typography } from '@mui/material';

const stats = [
  { label: 'Clientes ativos', value: '—' },
  { label: 'Agendamentos hoje', value: '—' },
  { label: 'Receita do dia', value: '—' },
];

const DashboardPage: React.FC = () => {
  return (
    <>
      <Typography variant="h4" component="h1" gutterBottom>
        Visão geral
      </Typography>
      <Grid container spacing={3}>
        {stats.map((stat) => (
          <Grid key={stat.label} xs={12} md={4}>
            <Card>
              <CardContent>
                <Typography color="text.secondary" gutterBottom>
                  {stat.label}
                </Typography>
                <Typography variant="h5">{stat.value}</Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </>
  );
};

export default DashboardPage;
