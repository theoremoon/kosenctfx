FROM ubuntu:18.04

LABEL maintainer "ptr-yudai <ptr.yudai@gmail.com>"

RUN apt update -y && apt upgrade -y
RUN apt install -y python3 socat binutils nasm build-essential
RUN adduser --disabled-password pwn

ADD challenge/run.sh    /etc/run.sh
RUN chmod 555 /etc/run.sh

WORKDIR /home/pwn/
ADD challenge/flag.txt  ./
ADD challenge/server.py ./
RUN chmod 440 flag.txt
RUN chmod 550 server.py

RUN chown -R root:pwn /home/pwn
RUN chmod 1733 /tmp /var/tmp /dev/shm

USER pwn
