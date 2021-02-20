#!/bin/bash

printf "\n\n### Restart database? [Y/N]"
read startdb
if [ "$startdb" == "Y" ]; then
    # Start MongoDB
    ( cd ../scripts/database; sh ./restart.sh )
fi

printf "\n\n### Reset admin user? [Y/N]"
read resetAdmin
if [ "$resetAdmin" == "Y" ]; then
    # Reset admin user
    ( cd ./admin/user; sh ./resetUser.sh )
    if [ $? -ne 0 ]; then
        exit 1
    fi
fi

printf "\n\n### Keystore ###\n"
printf "\n   Generate new certification chain, root to leaf? [Y/N]"
read pem
if [ "$pem" == "Y" ]; then
    ( cd ./keystore; sh ./generateRootCAKeys.sh )
    if [ $? -ne 0 ]; then
        exit 1
    fi
    ( cd ./keystore; sh ./generateServerCAKeys.sh )
    if [ $? -ne 0 ]; then
        exit 1
    fi
    ( cd ./keystore; sh ./generateServerKeys.sh )
    if [ $? -ne 0 ]; then
        exit 1
    fi
else
    printf "\n   Generate new Maths certification chain, issuer and leaf? [Y/N]"
    read pem
    if [ "$pem" == "Y" ]; then
        ( cd ./keystore; sh ./generateServerCAKeys.sh )
        if [ $? -ne 0 ]; then
            exit 1
        fi
        ( cd ./keystore; sh ./generateServerKeys.sh )
        if [ $? -ne 0 ]; then
            exit 1
        fi
    else
        printf "\n   Generate new Maths certificate? [Y/N]"
        read pem
        if [ "$pem" == "Y" ]; then
            ( cd ./keystore; sh ./generateServerKeys.sh )
            if [ $? -ne 0 ]; then
                exit 1
            fi
        fi
    fi
fi
