CREATE TABLE inventory_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES companies(id),
    product_id UUID NOT NULL REFERENCES products(id),
    order_id UUID REFERENCES sales_orders(id),
    type VARCHAR(16) NOT NULL,
    quantity INT NOT NULL CHECK (quantity >= 0),
    reason VARCHAR(160),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_inventory_tenant_product ON inventory_movements (tenant_id, product_id);

CREATE TRIGGER set_timestamp_inventory
BEFORE UPDATE ON inventory_movements
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE OR REPLACE FUNCTION apply_inventory_movement() RETURNS TRIGGER AS $$
DECLARE
    current_stock INT;
BEGIN
    SELECT stock_qty INTO current_stock FROM products WHERE id = NEW.product_id FOR UPDATE;
    IF current_stock IS NULL THEN
        RAISE EXCEPTION 'Produto % não encontrado para movimento de estoque', NEW.product_id;
    END IF;

    IF NEW.type = 'in' THEN
        current_stock := current_stock + NEW.quantity;
    ELSIF NEW.type = 'out' THEN
        IF current_stock < NEW.quantity THEN
            RAISE EXCEPTION 'Estoque insuficiente para produto %', NEW.product_id;
        END IF;
        current_stock := current_stock - NEW.quantity;
    ELSIF NEW.type = 'adjustment' THEN
        current_stock := NEW.quantity;
    ELSE
        RAISE EXCEPTION 'Tipo de movimento inválido: %', NEW.type;
    END IF;

    UPDATE products SET stock_qty = current_stock, updated_at = NOW() WHERE id = NEW.product_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_inventory_movements
AFTER INSERT ON inventory_movements
FOR EACH ROW
EXECUTE FUNCTION apply_inventory_movement();
