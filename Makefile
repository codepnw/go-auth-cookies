include cmd/dev.env

docker-up:
	docker compose --env-file cmd/dev.env up -d

goose-up:
	cd internal/db/schema && goose postgres ${DB_URL} up

