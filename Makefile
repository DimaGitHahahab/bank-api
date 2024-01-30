include .env
export

compose-up:
	docker compose up --build && docker compose logs --follow

compose-down:
	docker compose down --remove-orphans

migrate-create:
	migrate create -ext sql -dir migrations 'init_schema'

migrate-up:
	migrate -path migrations -database '$(DB_URL)?sslmode=disable' up