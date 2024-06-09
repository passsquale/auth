#!/bin/bash
source .env

sleep 2 && /bin/goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v