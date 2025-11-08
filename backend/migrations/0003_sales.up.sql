CREATE TABLE sales_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES companies(id),
    client_id UUID NOT NULL REFERENCES clients(id),
    booking_id UUID REFERENCES bookings(id),
    status VARCHAR(32) NOT NULL,
    payment_method VARCHAR(32),
    total NUMERIC(12,2) NOT NULL DEFAULT 0,
    discount NUMERIC(12,2) NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_sales_orders_tenant_status ON sales_orders (tenant_id, status);

CREATE TRIGGER set_timestamp_sales_orders
BEFORE UPDATE ON sales_orders
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE sales_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES companies(id),
    order_id UUID NOT NULL REFERENCES sales_orders(id) ON DELETE CASCADE,
    item_type VARCHAR(16) NOT NULL,
    item_ref_id UUID NOT NULL,
    quantity INT NOT NULL,
    unit_price NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_sales_items_order ON sales_items (order_id);

CREATE TRIGGER set_timestamp_sales_items
BEFORE UPDATE ON sales_items
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES companies(id),
    order_id UUID NOT NULL REFERENCES sales_orders(id),
    method VARCHAR(32) NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    paid_at TIMESTAMPTZ NOT NULL,
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_payments_tenant_method ON payments (tenant_id, method);

CREATE TRIGGER set_timestamp_payments
BEFORE UPDATE ON payments
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
