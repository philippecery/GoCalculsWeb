#!/bin/bash

printf "\n\n### Restart database? [Y/N]"
read startdb
if [ "$startdb" == "Y" ]; then
    # Start MongoDB
    ( cd ../scripts/database; sh ./restart.sh )
fi

printf "\n\n### Reset database? [Y/N]"
read resetDB
if [ "$resetDB" == "Y" ]; then
    # Reset database
    ( cd ./mongo; sh ./reset.sh )
    if [ $? -ne 0 ]; then
        exit 1
    fi
fi

printf "\n\n### Keystore ###\n"
printf "\n   Generate new certification chain, root to leaf? [Y/N]"
read all
if [ "$all" == "Y" ]; then
    ( cd ./keystore; sh ./generateAllKeys.sh )
    if [ $? -ne 0 ]; then
        exit 1
    fi
else
    sh ./setupService.sh database
    if [ $? -ne 0 ]; then
        exit 1
    fi
    sh ./setupService.sh webapp
    if [ $? -ne 0 ]; then
        exit 1
    fi
fi
