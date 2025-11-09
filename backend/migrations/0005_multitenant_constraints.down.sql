ALTER TABLE inventory_movements
    DROP CONSTRAINT IF EXISTS fk_inventory_movements_orders_tenant,
    DROP CONSTRAINT IF EXISTS fk_inventory_movements_products_tenant;

ALTER TABLE payments
    DROP CONSTRAINT IF EXISTS fk_payments_orders_tenant;

ALTER TABLE sales_items
    DROP CONSTRAINT IF EXISTS fk_sales_items_orders_tenant;

ALTER TABLE sales_orders
    DROP CONSTRAINT IF EXISTS fk_sales_orders_bookings_tenant,
    DROP CONSTRAINT IF EXISTS fk_sales_orders_clients_tenant;

ALTER TABLE bookings
    DROP CONSTRAINT IF EXISTS fk_bookings_services_tenant,
    DROP CONSTRAINT IF EXISTS fk_bookings_professionals_tenant,
    DROP CONSTRAINT IF EXISTS fk_bookings_clients_tenant;

ALTER TABLE sales_orders
    DROP CONSTRAINT IF EXISTS uq_sales_orders_tenant_id_id;

ALTER TABLE bookings
    DROP CONSTRAINT IF EXISTS uq_bookings_tenant_id_id;

ALTER TABLE products
    DROP CONSTRAINT IF EXISTS uq_products_tenant_id_id;

ALTER TABLE services
    DROP CONSTRAINT IF EXISTS uq_services_tenant_id_id;

ALTER TABLE professionals
    DROP CONSTRAINT IF EXISTS uq_professionals_tenant_id_id;

ALTER TABLE clients
    DROP CONSTRAINT IF EXISTS uq_clients_tenant_id_id;
