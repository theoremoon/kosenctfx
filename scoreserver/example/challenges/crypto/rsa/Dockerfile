FROM python:3.7-alpine

RUN apk --no-cache add socat
RUN apk add gcc g++ make libffi-dev openssl-dev
RUN pip3 install pycryptodome

COPY ./challenge/prob.py /prob.py

CMD socat TCP-L:11000,fork,reuseaddr EXEC:"python prob.py",
