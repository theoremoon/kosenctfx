help:
	:

up:
	(cd dev; docker-compose up)
down:
	(cd dev; docker-compose down)

build:
	mkdir bin
	(cd scoreserver; go build -o ../bin)
	(cd challengemanager; go build -o ../bin)

run: build
	(source ./envfile; cd scoreserver; ./scoreserver)

sql:
	(cd dev; docker-compose exec db mysql -u kosenctfxuser -pkosenctfxpassword kosenctfx)
