#!/bin/bash

# Load .env file
if [ -f .env ]; then
    export $(cat .env | sed 's/#.*//g' | xargs)
fi

# Run the migration command
migrate -path migrations -database "$DATABASE_URL" "$@"