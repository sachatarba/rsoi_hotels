# Строка подключения (для локального запуска)
DB_DSN := "postgres://postgres:secret@localhost:5432/reservation_db?sslmode=disable"

# Создать новую миграцию: make migration name=add_orders
migration:
	goose -dir migrations create $(name) sql

# Накатить миграции (локально)
migrate-up:
	goose -dir migrations postgres $(DB_DSN) up

# Откатить одну миграцию назад
migrate-down:
	goose -dir migrations postgres $(DB_DSN) down

# Посмотреть статус (какие миграции применены)
migrate-status:
	goose -dir migrations postgres $(DB_DSN) status