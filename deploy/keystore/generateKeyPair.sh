#!/bin/bash

SERVICE_NAME=$1
if [[ -z $SERVICE_NAME ]]; then
    echo "Usage: generateAllServiceKeys.sh SERVICE_NAME [UNPROTECTED]"
    echo "    where SERVICE_NAME is the name of the service to generate the keys for"
    echo "      and if there is a second parameter, the private won't be password-protected"
    exit 1
fi

UNPROTECTED=$2
if [[ -z $UNPROTECTED ]]; then
    UNPROTECTED=false
fi

OPENSSL_CMD=openssl
TARGET_PATH=../../config/tls/$SERVICE_NAME

if [ "$UNPROTECTED" == "false" ]; then
    printf "Enter the Maths $SERVICE_NAME private key password:"
    read -s password
    if [ ${#password} -ge 8 ]; then
        printf "\nEnter the Maths $SERVICE_NAME private key password again:"
        read -s confirmPassword
        if [ "$password" == "$confirmPassword" ]; then
            $OPENSSL_CMD genrsa 2048 | $OPENSSL_CMD pkcs8 -topk8 -out $TARGET_PATH/server.key -passout pass:$password -v2 aes256
            if [ $? -ne 0 ]; then
                cd ..
                echo "[ERROR] Maths $SERVICE_NAME key generation failed"
                exit 1
            fi
            $OPENSSL_CMD req -new -sha256 -key $TARGET_PATH/server.key -passin pass:$password -subj "/CN=maths-${SERVICE_NAME}" -out $TARGET_PATH/server.csr
            if [ $? -ne 0 ]; then
                cd ..
                echo "[ERROR] Maths $SERVICE_NAME CSR generation failed"
                exit 1
            fi
        else
            printf "\nERROR: You must enter the same private key password twice\n"
            exit 12
        fi
    else
        printf "\nERROR: The private key password must be at least 8 characters long\n"
        exit 13
    fi
else
    $OPENSSL_CMD genrsa 2048 | $OPENSSL_CMD pkcs8 -topk8 -nocrypt -out $TARGET_PATH/server.key
    if [ $? -ne 0 ]; then
        cd ..
        echo "[ERROR] Maths $SERVICE_NAME key generation failed"
        exit 1
    fi
    $OPENSSL_CMD req -new -sha256 -key $TARGET_PATH/server.key -subj "/CN=maths-${SERVICE_NAME}" -out $TARGET_PATH/server.csr
    if [ $? -ne 0 ]; then
        cd ..
        echo "[ERROR] Maths $SERVICE_NAME CSR generation failed"
        exit 1
    fi
fi