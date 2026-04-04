compose-build:
	docker compose build

compose-up:
	docker compose up -d

compose-up-proxy:
	docker compose --profile proxy up -d

compose-down:
	docker compose down

compose-logs:
	docker compose logs -f --tail=200

compose-config:
	docker compose config

server-test:
	cd Server && go test ./...

