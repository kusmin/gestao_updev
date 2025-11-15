import React from 'react';
import ResourceListPage from '@/components/ResourceListPage';
import SaleForm from './SaleForm';
import { SalesOrder } from '@/types/sales';

const columns = [
  { header: 'Cliente', accessor: 'client_id' as keyof SalesOrder },
  { header: 'Status', accessor: 'status' as keyof SalesOrder },
  { header: 'Total', accessor: 'total' as keyof SalesOrder },
];

const SaleListPage: React.FC = () => {
  return (
    <ResourceListPage<SalesOrder>
      title="Vendas"
      endpoint="/admin/sales/orders"
      columns={columns}
      renderForm={(props) => <SaleForm {...props} sale={props.item} />}
    />
  );
};

export default SaleListPage;
