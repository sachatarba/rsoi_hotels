-- +goose Up
-- +goose StatementBegin

INSERT INTO loyalty (id, username, reservation_count, status, discount)
VALUES (1, 'Test Max', 25, 'GOLD', 10)
    ON CONFLICT (username) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM loyalty WHERE username = 'Test Max';

-- +goose StatementEnd