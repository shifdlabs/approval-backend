up:
	docker compose up -d

down:
	docker compose down

seed:
	go run cmd/seed/main.go

# wipe everything and start fresh
fresh:
	docker compose down -v
	docker compose up -d
	sleep 3
	go run cmd/seed/main.go