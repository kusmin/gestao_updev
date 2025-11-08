CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE OR REPLACE FUNCTION trigger_set_timestamp() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(140) NOT NULL,
    document VARCHAR(32) UNIQUE,
    timezone VARCHAR(60) NOT NULL DEFAULT 'America/Sao_Paulo',
    phone VARCHAR(32),
    email VARCHAR(120),
    settings JSONB NOT NULL DEFAULT '{}'::jsonb,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TRIGGER set_timestamp_companies
BEFORE UPDATE ON companies
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(160) NOT NULL,
    email VARCHAR(160) NOT NULL,
    phone VARCHAR(32),
    role VARCHAR(32) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    profile JSONB NOT NULL DEFAULT '{}'::jsonb,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_users_email_tenant ON users (tenant_id, lower(email));

CREATE TRIGGER set_timestamp_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(160) NOT NULL,
    email VARCHAR(160),
    phone VARCHAR(32),
    notes TEXT,
    tags JSONB DEFAULT '[]'::jsonb,
    contact JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_clients_tenant_name ON clients (tenant_id, name);

CREATE TRIGGER set_timestamp_clients
BEFORE UPDATE ON clients
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(160) NOT NULL,
    category VARCHAR(80),
    description TEXT,
    duration_minutes INT NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    color VARCHAR(16),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_services_name_tenant ON services (tenant_id, lower(name));

CREATE TRIGGER set_timestamp_services
BEFORE UPDATE ON services
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES companies(id),
    name VARCHAR(160) NOT NULL,
    sku VARCHAR(80) NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    cost NUMERIC(12,2) DEFAULT 0,
    stock_qty INT NOT NULL DEFAULT 0,
    min_stock INT NOT NULL DEFAULT 0,
    description TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_products_sku_tenant ON products (tenant_id, lower(sku));

CREATE TRIGGER set_timestamp_products
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES companies(id),
    entity VARCHAR(64) NOT NULL,
    action VARCHAR(64) NOT NULL,
    actor_id UUID NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_tenant_entity ON audit_logs (tenant_id, entity);
