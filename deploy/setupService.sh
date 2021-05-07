#!/bin/bash

SERVICE_NAME=$1
if [[ -z $SERVICE_NAME ]]; then
    echo "This script will create the RSA key pairs of a new intermediary CA and the server for a specific service."
    echo "The intermediary CA certificate will be signed by the existing root CA and the server certificate will be signed by the newly created intermediary CA."
    echo "Usage: setupService.sh SERVICE_NAME"
    echo "    where SERVICE_NAME is the name of the service to setup, either webapp or database."
    exit 1
fi

printf "\n   Generate new Maths $SERVICE_NAME certification chain, issuer and leaf? [Y/N]"
read pem
if [ "$pem" == "Y" ]; then
    ( cd ./keystore; sh ./generateAllServiceKeys.sh $SERVICE_NAME)
    if [ $? -ne 0 ]; then
        exit 1
    fi
else
    printf "\n   Generate new Maths $SERVICE_NAME server certificate? [Y/N]"
    read pem
    if [ "$pem" == "Y" ]; then
        ( cd ./keystore; sh ./generateServiceKeys.sh $SERVICE_NAME)
        if [ $? -ne 0 ]; then
            exit 1
        fi
    fi
fi
