.SILENT:

run:
	CONFIG_PATH=./config/local.yaml go run ./cmd/avito_assignment/main.go

migrate:
	CONFIG_PATH=./config/local.yaml MIGRATIONS_PATH=file://./migrations/ go run ./cmd/migrator/main.go

tests:
	CONFIG_PATH=../../config/local.yaml go test ./... -v --cover      
