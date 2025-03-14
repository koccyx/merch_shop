.SILENT:

run:
	CONFIG_PATH=./config/local.yaml go run ./cmd/avito_assignment/main.go

docker-run:
	docker-compose up --build  

migrate:
	CONFIG_PATH=./config/local.yaml MIGRATIONS_PATH=file://./migrations/ go run ./cmd/migrator/main.go

test:
	CONFIG_PATH=../../../config/local.yaml go test ./... -v -cover
