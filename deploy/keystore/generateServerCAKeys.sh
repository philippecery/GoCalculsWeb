#!/bin/bash

OPENSSL_CMD=openssl
ROOT_PATH=../pem
TARGET_PATH=../pem

printf "Enter the Maths CA private key password:"
read -s password
if [ ${#password} -ge 8 ]; then
    printf "\nEnter the Maths CA private key password again:"
    read -s confirmPassword
    if [ "$password" = "$confirmPassword" ]; then
        $OPENSSL_CMD genrsa 2048 | $OPENSSL_CMD pkcs8 -topk8 -out $TARGET_PATH/CA_server.key -passout pass:$password -v2 aes256
        if [ $? -ne 0 ]; then
            cd ..
            echo "[ERROR] CA Maths key generation failed"
            exit 1
        fi
        $OPENSSL_CMD req -new -sha256 -key $TARGET_PATH/CA_server.key -passin pass:$password -subj "/CN=CA_SERVER" -out $TARGET_PATH/CA_server.csr
        if [ $? -ne 0 ]; then
            cd ..
            echo "[ERROR] CA Maths CSR generation failed"
            exit 2
        fi
        printf "Enter the Root CA private key password:"
        read -s password
        $OPENSSL_CMD x509 -req -in $TARGET_PATH/CA_server.csr -CA $ROOT_PATH/CA_root.crt -CAkey $ROOT_PATH/CA_root.key -passin pass:$password -CAcreateserial -sha256 -days 365 -extfile extServer.cnf -extensions CA_server_cert_ext -out $TARGET_PATH/CA_server.crt
        if [ $? -ne 0 ]; then
            cd ..
            echo "[ERROR] CA Maths cert generation failed"
            exit 3
        fi
        rm -f $TARGET_PATH/*.csr
        rm -f .srl
        cat $TARGET_PATH/server.key > $TARGET_PATH/server.pem
        cat $TARGET_PATH/server.crt >> $TARGET_PATH/server.pem
    else
        printf "\nERROR: You must enter the same Maths CA private key password twice\n"
        exit 4
    fi
else
    printf "\nERROR: The Maths CA private key password must be at least 8 characters long\n"
    exit 5
fi
