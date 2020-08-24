FROM alpine:3.12

RUN apk add --update --no-cache ca-certificates tzdata
ADD ./production/scoreserver /scoreserver
