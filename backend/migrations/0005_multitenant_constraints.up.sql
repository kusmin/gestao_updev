ALTER TABLE clients
    ADD CONSTRAINT uq_clients_tenant_id_id UNIQUE (tenant_id, id);

ALTER TABLE professionals
    ADD CONSTRAINT uq_professionals_tenant_id_id UNIQUE (tenant_id, id);

ALTER TABLE services
    ADD CONSTRAINT uq_services_tenant_id_id UNIQUE (tenant_id, id);

ALTER TABLE products
    ADD CONSTRAINT uq_products_tenant_id_id UNIQUE (tenant_id, id);

ALTER TABLE bookings
    ADD CONSTRAINT uq_bookings_tenant_id_id UNIQUE (tenant_id, id);

ALTER TABLE sales_orders
    ADD CONSTRAINT uq_sales_orders_tenant_id_id UNIQUE (tenant_id, id);

ALTER TABLE bookings
    ADD CONSTRAINT fk_bookings_clients_tenant FOREIGN KEY (tenant_id, client_id) REFERENCES clients(tenant_id, id),
    ADD CONSTRAINT fk_bookings_professionals_tenant FOREIGN KEY (tenant_id, professional_id) REFERENCES professionals(tenant_id, id),
    ADD CONSTRAINT fk_bookings_services_tenant FOREIGN KEY (tenant_id, service_id) REFERENCES services(tenant_id, id);

ALTER TABLE sales_orders
    ADD CONSTRAINT fk_sales_orders_clients_tenant FOREIGN KEY (tenant_id, client_id) REFERENCES clients(tenant_id, id),
    ADD CONSTRAINT fk_sales_orders_bookings_tenant FOREIGN KEY (tenant_id, booking_id) REFERENCES bookings(tenant_id, id);

ALTER TABLE sales_items
    ADD CONSTRAINT fk_sales_items_orders_tenant FOREIGN KEY (tenant_id, order_id) REFERENCES sales_orders(tenant_id, id);

ALTER TABLE payments
    ADD CONSTRAINT fk_payments_orders_tenant FOREIGN KEY (tenant_id, order_id) REFERENCES sales_orders(tenant_id, id);

ALTER TABLE inventory_movements
    ADD CONSTRAINT fk_inventory_movements_products_tenant FOREIGN KEY (tenant_id, product_id) REFERENCES products(tenant_id, id),
    ADD CONSTRAINT fk_inventory_movements_orders_tenant FOREIGN KEY (tenant_id, order_id) REFERENCES sales_orders(tenant_id, id);
