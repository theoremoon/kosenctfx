FROM golang:1.17 as builder

WORKDIR /
RUN wget https://github.com/orisano/wayt/releases/download/v0.3.1/wayt_0.3.1_Linux_x86_64.tar.gz && tar xf wayt_0.3.1_Linux_x86_64.tar.gz

WORKDIR /go/src/github.com/theoremoon/kosenctfx/scoreserver
COPY . .
RUN go build -a -tags netgo -installsuffix netgo -ldflags="-extldflags \"-static\""

FROM alpine:3.12
RUN apk add --update --no-cache ca-certificates tzdata
COPY --from=builder /wayt /wayt
COPY --from=builder /go/src/github.com/theoremoon/kosenctfx/scoreserver/scoreserver /scoreserver


CMD sh -c '/wayt sql -dsn "$DBDSN" && /scoreserver'

