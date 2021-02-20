#!/bin/bash

OPENSSL_CMD=openssl
CA_PATH=../pem
TARGET_PATH=../pem

printf "Enter the Maths private key password:"
read -s password
if [ ${#password} -ge 8 ]; then
    printf "\nEnter the Maths private key password again:"
    read -s confirmPassword
    if [ "$password" = "$confirmPassword" ]; then
        $OPENSSL_CMD genrsa 2048 | $OPENSSL_CMD pkcs8 -topk8 -out $TARGET_PATH/server.key -passout pass:$password -v2 aes256
        if [ $? -ne 0 ]; then
            cd ..
            echo "[ERROR] Maths key generation failed"
            exit 1
        fi
        $OPENSSL_CMD req -new -sha256 -key $TARGET_PATH/server.key -passin pass:$password -subj "/CN=localhost" -out $TARGET_PATH/server.csr
        if [ $? -ne 0 ]; then
            cd ..
            echo "[ERROR] Maths CSR generation failed"
            exit 1
        fi
        printf "Enter the Maths CA private key password:"
        read -s password
        $OPENSSL_CMD x509 -req -in $TARGET_PATH/server.csr -CA $CA_PATH/CA_server.crt -CAkey $CA_PATH/CA_server.key -passin pass:$password -CAcreateserial -sha256 -days 365 -extfile extServer.cnf -extensions server_cert_ext -out $TARGET_PATH/server.crt
        if [ $? -ne 0 ]; then
            cd ..
            echo "[ERROR] Maths cert generation failed"
            exit 1
        fi
        rm -f $TARGET_PATH/*.csr
        rm -f .srl
        cat $TARGET_PATH/server.key > $TARGET_PATH/server.pem
        cat $TARGET_PATH/server.crt >> $TARGET_PATH/server.pem
    else
        printf "\nERROR: You must enter the same Maths private key password twice\n"
        exit 12
    fi
else
    printf "\nERROR: The Maths private key password must be at least 8 characters long\n"
    exit 13
fi
