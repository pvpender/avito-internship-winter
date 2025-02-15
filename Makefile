up:
	docker compose up -d --build

down:
	docker compose down

clear:
	docker compose down -v --remove-orphans
	docker compose rm -vsf