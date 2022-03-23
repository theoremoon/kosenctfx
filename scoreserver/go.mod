module github.com/theoremoon/kosenctfx/scoreserver

go 1.14

require (
	entgo.io/ent v0.9.1
	github.com/aws/aws-sdk-go v1.30.29
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-resty/resty/v2 v2.5.0
	github.com/go-sql-driver/mysql v1.5.1-0.20200311113236-681ffa848bae
	github.com/google/uuid v1.3.0
	github.com/graph-gophers/dataloader/v6 v6.0.0
	github.com/labstack/echo/v4 v4.1.16
	github.com/labstack/gommon v0.3.0
	github.com/mattn/anko v0.1.8
	github.com/minio/minio-go/v7 v7.0.15
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
	github.com/pariz/gountries v0.0.0-20200430155801-1c6a393df9c7
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	golang.org/x/crypto v0.0.0-20201216223049-8b5274cf687f
	golang.org/x/mod v0.4.2
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.7
	syreclabs.com/go/faker v1.2.3
)
