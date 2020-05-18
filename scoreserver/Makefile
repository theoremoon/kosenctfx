up:
	(cd dev; docker-compose up)
down:
	(cd dev; docker-compose down)

setchallenge:
	(cd cmd/setchallenge/; go build -a -tags netgo -installsuffix netgo -ldflags="-extldflags \"-static\"" -o ../../bin/setchallenge)

set-challenge:
	(cd bin; DBDSN='kosenctfxuser:kosenctfxpassword@tcp(localhost:13306)/kosenctfx' ./setchallenge -path="../example/challenges" -transfersh "http://transfer:password@localhost:9999/")
