migrate-create:
	docker compose run --rm --entrypoint migrate migrate create -ext sql -dir /migrations $(name)

sqlc:
	docker compose run --rm sqlc