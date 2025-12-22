CREATE TABLE IF NOT EXISTS categories(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_archived BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS products(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) UNIQUE NOT NULL,
    price BIGINT NOT NULL,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_archived BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS product_embedding_version(
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL UNIQUE REFERENCES products(id) ON DELETE RESTRICT,
    embedding_version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    is_archived BOOLEAN DEFAULT false
);

-- outbox таблица, не нужна нормализованность для производительности
CREATE TABLE outbox_events (
    id BIGSERIAL PRIMARY KEY,
    event_id UUID NOT NULL UNIQUE,
    product_id BIGINT NOT NULL,
    event_type VARCHAR(50) NOT NULL, -- product_event
    payload BYTEA NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    processing_started_at TIMESTAMP,
    processed_at TIMESTAMP
);

CREATE INDEX idx_outbox_pending ON outbox_events(status, created_at);