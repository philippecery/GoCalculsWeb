#!/bin/bash

printf "Enter the admin email address: "
read emailAddress
adminId=$(cd ../../tools/protect/;go run protect.go ../../../config/config.json -userId $emailAddress)
adminEmail=$(cd ../../tools/protect/;go run protect.go ../../../config/config.json -pii $emailAddress)

printf "\nEnter the admin password:"
read -s password
if [ ${#password} -ge 8 ]; then
    printf "\nEnter the admin password again:"
    read -s confirmPassword
    if [ "$password" = "$confirmPassword" ]; then
        printf "\nAdmin user, if exists, will be reset. Enter YES to confirm, anything else to cancel: "
        read confirmReset
        if [ "$confirmReset" = "YES" ]; then
            adminPassword=$(cd ../../tools/protect/;go run protect.go ../../../config/config.json -password $password)
            mongo --eval 'var adminId = "'"$adminId"'"; var adminEmail = "'"$adminEmail"'"; var adminPassword = "'"$adminPassword"'";' resetUser.js
        else
            printf "User 'admin' reset was CANCELLED\n"
            exit 11
        fi
    else
        printf "\nERROR: You must enter the same admin password twice\n"
        exit 12
    fi
else
    printf "\nERROR: The admin password must be at least 8 characters long\n"
    exit 13
fi
