#!/bin/bash
sudo apt install mkcert

sudo apt install libnss3-tools

mkcert -install

mkcert localhost

go build -o server cmd/main.go

sudo ./server
