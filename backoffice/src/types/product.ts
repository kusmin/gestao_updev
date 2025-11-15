export interface Product {
  id: string;
  tenant_id: string;
  name: string;
  description?: string;
  price: number;
  cost?: number;
  sku: string;
  min_stock?: number;
  stock_qty?: number;
  created_at: string;
  updated_at: string;
}
