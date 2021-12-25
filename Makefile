.PHONY: help
help:
	@cat Makefile | grep -E "^[A-Za-z0-9-]+:"

up:
	docker-compose up --build

down:
	docker-compose down --remove-orphans

build-production:
	mkdir -p production
	(cd scoreserver; go build -o ../production -a -tags netgo -installsuffix netgo -ldflags="-extldflags \"-static\"")

build: generate
	mkdir -p bin
	(cd scoreserver; go build -o ../bin; go build -o ../bin ./cmd/...)

generate:
	(cd scoreserver; go generate ./...)

sql:
	docker-compose exec db mysql -u kosenctfxuser -pkosenctfxpassword kosenctfx

pass:
	echo 'select token from configs;' | docker-compose exec -T db mysql -u kosenctfxuser -pkosenctfxpassword kosenctfx

test:
	docker-compose -f docker-compose.test.yml run go-test
