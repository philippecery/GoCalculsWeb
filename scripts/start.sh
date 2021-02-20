#!/bin/bash

( cd ./database ; sh ./restart.sh )
( cd ../webapp ; go run main.go ../config/config.json )