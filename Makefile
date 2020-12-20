help:
	:

up:
	parallel --lb sh -c ::: 'docker-compose up --build' 'cd ui; yarn watch'

down:
	docker-compose down

build-production:
	mkdir -p production
	(cd scoreserver; go build -o ../production -a -tags netgo -installsuffix netgo -ldflags="-extldflags \"-static\"")

build: generate
	mkdir -p bin
	(cd scoreserver; go build -o ../bin)

build-ui:
	(cd ui; yarn build)

generate:
	(cd scoreserver; go generate ./...)

run: build
	(source ./envfile; ./bin/scoreserver)

sql:
	docker-compose exec db mysql -u kosenctfxuser -pkosenctfxpassword kosenctfx

pass:
	echo 'select token from configs;' | docker-compose exec -T db mysql -u kosenctfxuser -pkosenctfxpassword kosenctfx
