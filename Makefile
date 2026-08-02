include .env
export
export POSTGRES_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable

postgres-up:
	@docker compose up postgres -d

postgres-down:
	@docker compose down postgres

env-clean:
	@rm -rf ./out

migrate-create:
	@docker compose run --rm migrate-postgres create -ext sql -dir ./migrations -seq $(seq)

migrate-up:
	@docker compose run --rm migrate-postgres -path ./migrations -database "$(POSTGRES_URL)" up

proto:
	@protoc --proto_path=api/proto \
	       --go_out=pkg/pb --go_opt=paths=source_relative \
	       --go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
	       api/proto/*.proto

