-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS loyalty
(
    id                SERIAL PRIMARY KEY,
    username          VARCHAR(80) NOT NULL UNIQUE,
    reservation_count INT         NOT NULL DEFAULT 0,
    status            VARCHAR(80) NOT NULL DEFAULT 'BRONZE'
    CHECK (status IN ('BRONZE', 'SILVER', 'GOLD')),
    discount          INT         NOT NULL
    );

COMMENT ON COLUMN loyalty.status IS 'Status: BRONZE, SILVER, GOLD';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS loyalty;

-- +goose StatementEnd