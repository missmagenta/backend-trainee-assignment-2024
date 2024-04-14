migrate-create:
	go run ./cmd/migrate db create_sql "$(name)"
.PHONY: migrate-create

migrate-run:
	go run ./cmd/migrate db migrate
.PHONY: migrate-run

migrate-rollback:
	go run ./cmd/migrate db rollback
.PHONY: migrate-rollback

integration-test:
	CONFIG_PATH=./../config/config.yml go test -v ./integration-test/...
.PHONY: integration-test