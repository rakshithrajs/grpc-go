include .env

UMS_DIR=UMS
FILES_DIR=files

UMS_MIGRATE_DB_URL="postgres://${UMS_DB_USER}:${UMS_DB_PASSWORD}@${UMS_DB_HOST}:${UMS_DB_PORT}/${UMS_DB_NAME}?sslmode=${UMS_DB_SSLMODE}"
FILES_MIGRATE_DB_URL="postgres://${FILES_DB_USER}:${FILES_DB_PASSWORD}@${FILES_DB_HOST}:${FILES_DB_PORT}/${FILES_DB_NAME}?sslmode=${FILES_DB_SSLMODE}"

check-UMS:
	@echo ${UMS_MIGRATE_DB_URL}

check-files:
	@echo ${FILES_MIGRATE_DB_URL}

create-UMS-migration:
	@migrate create -ext sql -dir $(UMS_DIR)/internal/storage/migrations -seq $(NAME)

create-files-migration:
	@migrate create -ext sql -dir $(FILES_DIR)/internal/storage/migrations -seq $(NAME)

migrate-UMS-up:
	@migrate -path $(UMS_DIR)/internal/storage/migrations -database $(UMS_MIGRATE_DB_URL) -verbose up

migrate-files-up:
	@migrate -path $(FILES_DIR)/internal/storage/migrations -database $(FILES_MIGRATE_DB_URL) -verbose up

migrate-UMS-down:
	@migrate -path $(UMS_DIR)/internal/storage/migrations -database $(UMS_MIGRATE_DB_URL) -verbose down 1

migrate-files-down:
	@migrate -path $(FILES_DIR)/internal/storage/migrations -database $(FILES_MIGRATE_DB_URL) -verbose down 1

UMS-version:
	@migrate -path ${UMS_DIR}/internal/storage/migrations -database ${UMS_MIGRATE_DB_URL} version

files-version:
	@migrate -path ${FILES_DIR}internal/storage/migrations -database ${FILES_MIGRATE_DB_URL} version

force-UMS: 
	@migrate -path ${UMS_DIR}/internal/storage/migrations -database ${UMS_MIGRATE_DB_URL} force ${VERSION}

force-files: 
	@migrate -path ${FILES_DIR}/internal/storage/migrations -database ${FILES_MIGRATE_DB_URL} force ${VERSION}

.PHONY: check-UMS check-files migrate-UMS-up migrate-files-up migrate-UMS-down migrate-files-down create-UMS-migration create-files-migration UMS-version files-version force-UMS force-files