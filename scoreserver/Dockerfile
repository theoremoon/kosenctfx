FROM golang:1.17 as builder

WORKDIR /
RUN wget https://github.com/orisano/wayt/releases/download/v0.3.1/wayt_0.3.1_Linux_x86_64.tar.gz && tar xf wayt_0.3.1_Linux_x86_64.tar.gz

WORKDIR /go/src/github.com/theoremoon/kosenctfx/scoreserver
COPY . .
RUN go build -o /scoreserver -a -tags netgo -installsuffix netgo -ldflags="-extldflags \"-static\"" ./cmd/kosenctfx 

FROM debian:12-slim
ENV DEBIAN_FRONTEND=noninteractive
COPY --from=builder /wayt /wayt
COPY --from=builder /scoreserver /scoreserver
RUN apt update && apt install -yqq ca-certificates


CMD sh -c '/wayt sql -dsn "$DBDSN" && /scoreserver'

