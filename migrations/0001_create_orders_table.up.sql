CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_id uuid NOT NULL UNIQUE,
    order_status text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);