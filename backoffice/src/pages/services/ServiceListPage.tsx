import React from 'react';
import ResourceListPage from '@/components/ResourceListPage';
import ServiceForm from './ServiceForm';
import { Service } from '@/types/service';

const columns = [
  { header: 'Nome', accessor: 'name' as keyof Service },
  { header: 'Duração (min)', accessor: 'duration_minutes' as keyof Service },
  { header: 'Preço', accessor: 'price' as keyof Service },
];

const ServiceListPage: React.FC = () => {
  return (
    <ResourceListPage<Service>
      title="Services"
      endpoint="/admin/services"
      columns={columns}
      renderForm={(props) => <ServiceForm {...props} service={props.item} />}
    />
  );
};

export default ServiceListPage;
