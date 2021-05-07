#!/bin/bash

SERVICE_NAME=$1
if [[ -z $SERVICE_NAME ]]; then
    echo "Usage: generateServiceKeys.sh SERVICE_NAME"
    echo "    where SERVICE_NAME is the name of the service to generate the keys for"
    exit 1
fi

OPENSSL_CMD=openssl
CA_PATH=../../config/tls/$SERVICE_NAME
TARGET_PATH=../../config/tls/$SERVICE_NAME

MERGE=true
if [ "$SERVICE_NAME" == "database" ]; then
    UNPROTECTED=true
    MERGE=false
fi

sh ./generateKeyPair.sh $SERVICE_NAME $UNPROTECTED
if [ $? -ne 0 ]; then
    cd ..
    echo "[ERROR] Maths $SERVICE_NAME key pair or CSR generation failed"
    exit 1
fi

printf "Enter the Maths $SERVICE_NAME CA private key password:"
read -s password
$OPENSSL_CMD x509 -req -in $TARGET_PATH/server.csr -CA $CA_PATH/CA_server.crt -CAkey $CA_PATH/CA_server.key -passin pass:$password -CAcreateserial -sha256 -days 365 -extfile extServer.cnf -extensions ${SERVICE_NAME}_cert_ext -out $TARGET_PATH/server.crt
if [ $? -ne 0 ]; then
    cd ..
    echo "[ERROR] Maths $SERVICE_NAME cert generation failed"
    exit 1
fi

rm -f $TARGET_PATH/*.csr
rm -f .srl

if [ "$MERGE" == "true" ]; then
    cat $TARGET_PATH/server.key > $TARGET_PATH/server.pem
    cat $TARGET_PATH/server.crt >> $TARGET_PATH/server.pem
    rm -f $TARGET_PATH/server.key
    rm -f $TARGET_PATH/server.crt
fi
