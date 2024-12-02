#!/usr/bin/env bash

make clean && make linux-amd64

scp -r bin/ss-proxy-linux-amd64 root@chuangjie.icu:/data/clash-ws/clash-ws
