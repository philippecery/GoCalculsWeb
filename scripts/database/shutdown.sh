#!/bin/bash
if [ -f ~/mongod.pid ]; then
    PID=$(cat ~/mongod.pid)
    while $(kill $PID 2>/dev/null); do
        sleep 1
    done
    rm ~/mongod.pid
    printf "MongoDB should be down...\n"
fi