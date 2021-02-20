#!/bin/bash

OPENSSL_CMD=openssl
TARGET_PATH=../pem

printf "Enter the Root CA private key password:"
read -s password
if [ ${#password} -ge 8 ]; then
    printf "\nEnter the Root CA private key password again:"
    read -s confirmPassword
    if [ "$password" = "$confirmPassword" ]; then
        $OPENSSL_CMD genrsa 2048 | $OPENSSL_CMD pkcs8 -topk8 -out $TARGET_PATH/CA_root.key -passout pass:$password -v2 aes256
        if [ $? -ne 0 ]; then
            cd ..
            echo "[ERROR] CA Root key generation failed"
            exit 1
        fi
        $OPENSSL_CMD req -x509 -new -key $TARGET_PATH/CA_root.key -passin pass:$password -sha256 -days 365 -subj "/CN=CA_ROOT" -config extRoot.cnf -out $TARGET_PATH/CA_root.crt
        if [ $? -ne 0 ]; then
            cd ..
            echo "[ERROR] CA Root cert generation failed"
            exit 1
        fi
    else
        printf "\nERROR: You must enter the same Root CA private key password twice\n"
        exit 12
    fi
else
    printf "\nERROR: The Root CA private key password must be at least 8 characters long\n"
    exit 13
fi

rm -f .srl