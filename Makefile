help:
	:

up:
	(cd dev; docker-compose up)
down:
	(cd dev; docker-compose down)

.PHONY: bucket
bucket:
	(cd bucket; docker-compose up)

bucket-down:
	(cd bucket; docker-compose down)

setchallenge:
	(cd scoreserver/cmd/setchallenge/; go build -a -tags netgo -installsuffix netgo -ldflags="-extldflags \"-static\"" -o ../../../bin/setchallenge)

set-challenge:
	(cd bin; DBDSN='kosenctfxuser:kosenctfxpassword@tcp(localhost:13306)/kosenctfx' ./setchallenge -path="../example/challenges" -transfersh "http://transfer:password@localhost:9999/")

build:
	(cd scoreserver; go build)

run: build
	(source ./envfile; cd scoreserver; ./scoreserver)

sql:
	(cd dev; docker-compose exec db mysql -u kosenctfxuser -pkosenctfxpassword kosenctfx)
