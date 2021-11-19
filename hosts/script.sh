#!/bin/sh

txeh add 127.0.0.1 minio
trap 'txeh remove host minio; exit' INT HUP TERM

while :; do
    sleep 1
done
