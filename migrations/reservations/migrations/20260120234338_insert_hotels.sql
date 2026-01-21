-- +goose Up
-- +goose StatementBegin

INSERT INTO hotels (id, hotel_uid, name, country, city, address, stars, price)
VALUES (1, '049161bb-badd-4fa8-9d90-87c9a82b0668', 'Ararat Park Hyatt Moscow', 'Россия', 'Москва', 'Неглинная ул., 4', 5, 10000)
    ON CONFLICT (id) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM hotels WHERE id = 1;

-- +goose StatementEnd