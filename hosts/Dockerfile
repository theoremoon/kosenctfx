FROM ubuntu:latest

RUN apt-get update && apt-get install -y curl

ADD ./script.sh /script.sh

WORKDIR /opt
RUN curl -fsSL -O https://github.com/txn2/txeh/releases/download/v1.3.0/txeh_linux_amd64.tar.gz && tar xf txeh_linux_amd64.tar.gz && mv txeh /usr/bin/txeh

WORKDIR /

