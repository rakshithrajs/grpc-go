include .env

UMS_DIR=UMS
MMS_DIR=MMS

UMS_MIGRATE_DB_URL="postgres://${UMS_DB_USER}:${UMS_DB_PASSWORD}@${UMS_DB_HOST}:${UMS_DB_PORT}/${UMS_DB_NAME}?sslmode=${UMS_DB_SSLMODE}"
MMS_MIGRATE_DB_URL="postgres://${MMS_DB_USER}:${MMS_DB_PASSWORD}@${MMS_DB_HOST}:${MMS_DB_PORT}/${MMS_DB_NAME}?sslmode=${MMS_DB_SSLMODE}"

check-UMS:
	@echo ${UMS_MIGRATE_DB_URL}

check-MMS:
	@echo ${MMS_MIGRATE_DB_URL}

create-UMS-migration:
	@migrate create -ext sql -dir $(UMS_DIR)/internal/storage/migrations -seq $(NAME)

create-MMS-migration:
	@migrate create -ext sql -dir $(MMS_DIR)/internal/storage/migrations -seq $(NAME)

migrate-UMS-up:
	@migrate -path $(UMS_DIR)/internal/storage/migrations -database $(UMS_MIGRATE_DB_URL) -verbose up

migrate-MMS-up:
	@migrate -path $(MMS_DIR)/internal/storage/migrations -database $(MMS_MIGRATE_DB_URL) -verbose up

migrate-UMS-down:
	@migrate -path $(UMS_DIR)/internal/storage/migrations -database $(UMS_MIGRATE_DB_URL) -verbose down 1

migrate-MMS-down:
	@migrate -path $(MMS_DIR)/internal/storage/migrations -database $(MMS_MIGRATE_DB_URL) -verbose down 1

UMS-version:
	@migrate -path ${UMS_DIR}/internal/storage/migrations -database ${UMS_MIGRATE_DB_URL} version

MMS-version:
	@migrate -path ${MMS_DIR}internal/storage/migrations -database ${MMS_MIGRATE_DB_URL} version

force-UMS: 
	@migrate -path ${UMS_DIR}/internal/storage/migrations -database ${UMS_MIGRATE_DB_URL} force ${VERSION}

force-MMS: 
	@migrate -path ${MMS_DIR}/internal/storage/migrations -database ${MMS_MIGRATE_DB_URL} force ${VERSION}

.PHONY: check-UMS check-MMS migrate-UMS-up migrate-MMS-up migrate-UMS-down migrate-MMS-down create-UMS-migration create-MMS-migration UMS-version MMS-version force-UMS force-MMS