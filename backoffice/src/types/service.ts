export interface Service {
  id: string;
  tenant_id: string;
  name: string;
  description?: string;
  price: number;
  duration_minutes: number;
  category?: string;
  color?: string;
  created_at: string;
  updated_at: string;
}
