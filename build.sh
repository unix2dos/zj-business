#!/bin/sh

GOOS=linux GOARCH=amd64 go build -v
rsync -avzP -e "ssh -i ~/.ssh/vultr" zj-business root@207.246.80.69:/opt/liuwei
