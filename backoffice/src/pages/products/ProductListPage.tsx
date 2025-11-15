import React from 'react';
import ResourceListPage from '@/components/ResourceListPage';
import ProductForm from './ProductForm';
import { Product } from '@/types/product';

const columns = [
  { header: 'Nome', accessor: 'name' as keyof Product },
  { header: 'SKU', accessor: 'sku' as keyof Product },
  { header: 'Preço', accessor: 'price' as keyof Product },
  { header: 'Estoque', accessor: 'stock_qty' as keyof Product },
];

const ProductListPage: React.FC = () => {
  return (
    <ResourceListPage<Product>
      title="Products"
      endpoint="/admin/products"
      columns={columns}
      renderForm={(props) => <ProductForm {...props} product={props.item} />}
    />
  );
};

export default ProductListPage;
