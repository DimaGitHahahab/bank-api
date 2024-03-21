compose-up:
	docker compose up --build && docker compose logs --follow

compose-down:
	docker compose down --remove-orphans

migrate-create:
	migrate create -ext sql -dir migrations 'init_schema'