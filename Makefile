include .env

ACCOUNT_DIR=account
FILES_DIR=files

ACCOUNT_MIGRATE_DB_URL="postgres://${ACCOUNT_DB_USER}:${ACCOUNT_DB_PASSWORD}@${ACCOUNT_DB_HOST}:${ACCOUNT_DB_PORT}/${ACCOUNT_DB_NAME}?sslmode=${ACCOUNT_DB_SSLMODE}"
FILES_MIGRATE_DB_URL="postgres://${FILES_DB_USER}:${FILES_DB_PASSWORD}@${FILES_DB_HOST}:${FILES_DB_PORT}/${FILES_DB_NAME}?sslmode=${FILES_DB_SSLMODE}"

check-account:
	@echo ${ACCOUNT_MIGRATE_DB_URL}

check-files:
	@echo ${FILES_MIGRATE_DB_URL}

create-account-migration:
	@migrate create -ext sql -dir $(ACCOUNT_DIR)/internal/storage/migrations -seq $(NAME)

create-files-migration:
	@migrate create -ext sql -dir $(FILES_DIR)/internal/storage/migrations -seq $(NAME)

migrate-account-up:
	@migrate -path $(ACCOUNT_DIR)/internal/storage/migrations -database $(ACCOUNT_MIGRATE_DB_URL) -verbose up

migrate-files-up:
	@migrate -path $(FILES_DIR)/internal/storage/migrations -database $(FILES_MIGRATE_DB_URL) -verbose up

migrate-account-down:
	@migrate -path $(ACCOUNT_DIR)/internal/storage/migrations -database $(ACCOUNT_MIGRATE_DB_URL) -verbose down 1

migrate-files-down:
	@migrate -path $(FILES_DIR)/internal/storage/migrations -database $(FILES_MIGRATE_DB_URL) -verbose down 1

account-version:
	@migrate -path ${ACCOUNT_DIR}/internal/storage/migrations -database ${ACCOUNT_MIGRATE_DB_URL} version

files-version:
	@migrate -path ${FILES_DIR}internal/storage/migrations -database ${FILES_MIGRATE_DB_URL} version

force-account: 
	@migrate -path ${ACCOUNT_DIR}/internal/storage/migrations -database ${ACCOUNT_MIGRATE_DB_URL} force ${VERSION}

force-files: 
	@migrate -path ${FILES_DIR}/internal/storage/migrations -database ${FILES_MIGRATE_DB_URL} force ${VERSION}

.PHONY: check-account check-files migrate-account-up migrate-files-up migrate-account-down migrate-files-down create-account-migration create-files-migration account-version files-version force-account force-files