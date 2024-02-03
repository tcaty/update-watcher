#!/bin/sh

# we should run this command from sh script
# cause only in this case we can use $POSTGRES_* env variables
migrate \
     "-database" \
     "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable" \
     "-path" \
     "/migrations" \
     "up"