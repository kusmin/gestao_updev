ALTER TABLE sales_orders
    ADD COLUMN IF NOT EXISTS payment_type VARCHAR(32);
