#!/bin/bash

SERVICE_NAME=$1
if [[ -z $SERVICE_NAME ]]; then
    echo "Usage: generateAllServiceKeys.sh SERVICE_NAME"
    echo "    where SERVICE_NAME is the name of the service to generate the keys for"
    exit 1
fi

OPENSSL_CMD=openssl
ROOT_PATH=../../config/tls
TARGET_PATH=../../config/tls/$SERVICE_NAME

printf "Enter the Maths $SERVICE_NAME CA private key password:"
read -s password
if [ ${#password} -ge 8 ]; then
    printf "\nEnter the Maths $SERVICE_NAME CA private key password again:"
    read -s confirmPassword
    if [ "$password" == "$confirmPassword" ]; then
        $OPENSSL_CMD genrsa 2048 | $OPENSSL_CMD pkcs8 -topk8 -out $TARGET_PATH/CA_server.key -passout pass:$password -v2 aes256
        if [ $? -ne 0 ]; then
            cd ..
            echo "[ERROR] Maths $SERVICE_NAME CA key generation failed"
            exit 1
        fi
        $OPENSSL_CMD req -new -sha256 -key $TARGET_PATH/CA_server.key -passin pass:$password -subj "/CN=ca_${SERVICE_NAME}_server" -out $TARGET_PATH/CA_server.csr
        if [ $? -ne 0 ]; then
            cd ..
            echo "[ERROR] Maths $SERVICE_NAME CA CSR generation failed"
            exit 2
        fi
        printf "Enter the Root CA private key password:"
        read -s password
        $OPENSSL_CMD x509 -req -in $TARGET_PATH/CA_server.csr -CA $ROOT_PATH/CA_root.crt -CAkey $ROOT_PATH/CA_root.key -passin pass:$password -CAcreateserial -sha256 -days 365 -extfile extServer.cnf -extensions CA_server_cert_ext -out $TARGET_PATH/CA_server.crt
        if [ $? -ne 0 ]; then
            cd ..
            echo "[ERROR] Maths $SERVICE_NAME CA cert generation failed"
            exit 3
        fi

        cat $ROOT_PATH/CA_root.crt > $TARGET_PATH/cacerts.pem
        cat $TARGET_PATH/CA_server.crt >> $TARGET_PATH/cacerts.pem

        sh ./generateServiceKeys.sh $SERVICE_NAME
        if [ $? -ne 0 ]; then
            exit 1
        fi

        rm -f $TARGET_PATH/*.csr
        rm -f .srl
    else
        printf "\nERROR: You must enter the same CA private key password twice\n"
        exit 4
    fi
else
    printf "\nERROR: The CA private key password must be at least 8 characters long\n"
    exit 5
fi
