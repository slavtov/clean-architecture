#!/bin/bash

if [ $TABLE ]; then
    migrate create \
        -ext sql \
        -dir $MIGRATE_PATH \
        -seq $TABLE;
else
    echo "var TABLE is missing";
fi
