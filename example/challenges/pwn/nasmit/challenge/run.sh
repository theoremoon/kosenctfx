#!/bin/sh
socat TCP-L:9001,reuseaddr,fork EXEC:/home/pwn/server.py
