include .env

MIGRATE_DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

check:
	@echo ${MIGRATE_DB_URL}

create:
	@migrate create -ext sql -dir internal/storage/migrations -seq ${NAME}

up:
	@migrate -path internal/storage/migrations -database ${MIGRATE_DB_URL} -verbose up

down: 
	@migrate -path internal/storage/migrations -database ${MIGRATE_DB_URL} -verbose down

version:
	@migrate -path internal/storage/migrations -database ${MIGRATE_DB_URL} version

force:
	@migrate -path internal/storage/migrations -database ${MIGRATE_DB_URL} force ${VERSION}

.PHONY: check create up down version force

