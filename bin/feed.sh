#!/bin/bash
LANG=C

cd $HOME/go/src/github.com/ohnishi/feed
git pull
sleep 10
go run github.com/ohnishi/feed/cmd reset
go run github.com/ohnishi/feed/cmd feed

date
