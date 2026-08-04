CREATE TABLE IF NOT EXISTS orders (
    id          BIGSERIAL                PRIMARY KEY,
    user_id     BIGINT                   NOT NULL,
    amount      NUMERIC(12, 2)           NOT NULL,
    status      VARCHAR(32)              NOT NULL DEFAULT 'CREATED',
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
