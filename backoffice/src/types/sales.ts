export interface SalesOrderItem {
  id: string;
  type: string;
  ref_id: string;
  quantity: number;
  unit_price: number;
}

export interface SalesOrder {
  id: string;
  tenant_id: string;
  client_id: string;
  booking_id?: string;
  discount?: number;
  notes?: string;
  items: SalesOrderItem[];
  created_at: string;
  updated_at: string;
}
