help:
	:

up:
	parallel --lb sh -c ::: 'docker-compose up --build' 'cd ui; yarn watch'

down:
	docker-compose down --remove-orphans

build-production:
	mkdir -p production
	(cd scoreserver; go build -o ../production -a -tags netgo -installsuffix netgo -ldflags="-extldflags \"-static\"")

build: generate
	mkdir -p bin
	(cd scoreserver; go build -o ../bin; go build -o ../bin ./cmd/...)

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

SIZE=20
seed:
	(docker-compose exec scoreserver go run cmd/seeder/main.go -size $(SIZE) -all)

seed-challenges:
	(docker-compose exec scoreserver go run cmd/seeder/main.go -size $(SIZE) -challenge)

seed-submissions:
	(docker-compose exec scoreserver go run cmd/seeder/main.go -size $(SIZE) -submission)

seed-teams:
	(docker-compose exec scoreserver go run cmd/seeder/main.go -size $(SIZE) -team)

