help:
	:

up:
	(cd dev; docker-compose up)

down:
	(cd dev; docker-compose down)

build:
	mkdir -p bin
	(cd scoreserver; go build -o ../bin)
	(cd challengemanager; go build -o ../bin)

run: build
	(source ./envfile; ./bin/scoreserver)

sql:
	(cd dev; docker-compose exec db mysql -u kosenctfxuser -pkosenctfxpassword kosenctfx)

