UID := $(shell id -u)
GID := $(shell id -g)

.PHONY: help
help:
	@cat Makefile | grep -E "^[A-Za-z0-9-]+:"

up:
	env UID=$(UID) GID=$(GID) docker-compose up --build

down:
	docker-compose down --remove-orphans

sql:
	docker-compose exec db mysql -u kosenctfxuser -pkosenctfxpassword kosenctfx

pass:
	echo 'select token from configs;' | docker-compose exec -T db mysql -u kosenctfxuser -pkosenctfxpassword kosenctfx

test:
	(cd scoreserver; go mod tidy)
	env UID=$(UID) GID=$(GID) docker compose -f compose.test.yaml down --remove-orphans
	env UID=$(UID) GID=$(GID) docker compose -f compose.test.yaml run --rm go-test
