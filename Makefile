help:
	:

up:
	(cd dev; docker-compose up)
down:
	(cd dev; docker-compose down)

setchallenge:
	(cd scoreserver/cmd/setchallenge/; go build -a -tags netgo -installsuffix netgo -ldflags="-extldflags \"-static\"" -o ../../../bin/setchallenge)

set-challenge:
	(cd bin; DBDSN='kosenctfxuser:kosenctfxpassword@tcp(localhost:13306)/kosenctfx' ./setchallenge -path="../example/challenges" -transfersh "http://transfer:password@localhost:9999/")

build:
	(cd scoreserver; go build)

run: build
	(cd scoreserver; DBDSN='kosenctfxuser:kosenctfxpassword@tcp(localhost:13306)/kosenctfx' FRONT='http://front.web.localhost:8080' ./scoreserver)

sql:
	(cd dev; docker-compose exec db mysql -u kosenctfxuser -pkosenctfxpassword kosenctfx)
