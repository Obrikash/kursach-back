include .envrc

# ========================================================================= #
# DEVELOPMENT
# ========================================================================= #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@go run ./cmd/api -db-dsn=${SWIMMING_POOL_DSN}

.PHONY: db/psql
db/psql:
	@psql ${SWIMMING_POOL_DSN}