help:
	:

up:
	docker-compose up

down:
	docker-compose down

build-production:
	mkdir -p production
	(cd scoreserver; go build -o ../production -a -tags netgo -installsuffix netgo -ldflags="-extldflags \"-static\"")

build:
	mkdir -p bin
	(cd scoreserver; go build -o ../bin)

build-ui:
	(cd ui; yarn build)

run: build
	(source ./envfile; ./bin/scoreserver)

sql:
	(cd dev; docker-compose exec db mysql -u kosenctfxuser -pkosenctfxpassword kosenctfx)

