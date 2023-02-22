#!/bin/bash

对应的编译语句
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./release/xxx.mac
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./release/xxx.exe
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o ./release/xxx.linux

chmod +x xxx.mac
touch nohup.out
nohup xxx.mac & tail -f nohup.out
ps -ef | grep xxx.mac | grep -v grep