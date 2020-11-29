#!/bin/bash

# Reset admin user
printf "Enter the admin password:"
read -s password
if [ ${#password} -ge 8 ]; then
    printf "\nEnter the admin password again:"
    read -s confirmPassword
    if [ "$password" = "$confirmPassword" ]; then
        printf "\nUser 'admin', if existing, will be reset. Enter YES to confirm, anything else to cancel: "
        read confirmReset
        if [ "$confirmReset" = "YES" ]; then
            hashPassword=$(go run hashPassword.go $password)
            mongo --eval 'var adminPassword = "'"$hashPassword"'"; var uri = "localhost:27017";' resetUser.js
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
