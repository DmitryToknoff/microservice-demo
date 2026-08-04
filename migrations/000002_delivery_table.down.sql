CREATE TABLE IF NOT EXISTS deliveries (
    id           BIGSERIAL       PRIMARY KEY,
    order_id     BIGINT          NOT NULL UNIQUE,
    address      VARCHAR(255)    NOT NULL DEFAULT 'Default',
    status       VARCHAR(32)     NOT NULL DEFAULT 'PROCESSING',
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

CREATE INDEX idx_deliveries_order_id ON deliveries(order_id);
