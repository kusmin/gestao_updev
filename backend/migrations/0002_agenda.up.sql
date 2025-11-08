CREATE TABLE professionals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES companies(id),
    user_id UUID REFERENCES users(id),
    name VARCHAR(160) NOT NULL,
    specialties JSONB NOT NULL DEFAULT '[]'::jsonb,
    max_parallel INT NOT NULL DEFAULT 1,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_professionals_tenant ON professionals (tenant_id, active);

CREATE TRIGGER set_timestamp_professionals
BEFORE UPDATE ON professionals
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE availability_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES companies(id),
    professional_id UUID NOT NULL REFERENCES professionals(id) ON DELETE CASCADE,
    weekday SMALLINT NOT NULL CHECK (weekday BETWEEN 0 AND 6),
    start_time VARCHAR(8) NOT NULL,
    end_time VARCHAR(8) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_availability_professional ON availability_rules (professional_id, weekday);

CREATE TRIGGER set_timestamp_availability
BEFORE UPDATE ON availability_rules
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES companies(id),
    client_id UUID NOT NULL REFERENCES clients(id),
    professional_id UUID NOT NULL REFERENCES professionals(id),
    service_id UUID NOT NULL REFERENCES services(id),
    status VARCHAR(32) NOT NULL,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    notes TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_bookings_conflict ON bookings (tenant_id, professional_id, start_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_bookings_status ON bookings (tenant_id, status);

CREATE TRIGGER set_timestamp_bookings
BEFORE UPDATE ON bookings
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
