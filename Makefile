DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
MIGRATIONS_DIR=migrations

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down
