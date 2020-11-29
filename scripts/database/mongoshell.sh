#!/bin/bash

printf "Enter the superuser password:"
read -s password
~/mongodb/bin/mongo --port 27017 -u superuser -p $password --authenticationDatabase 'admin'