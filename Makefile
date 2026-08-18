migrate-create:
	docker compose run --rm migrate $(name)

sqlc:
	docker compose run --rm sqlc