.PHONY: run

run:
	go run ./cmd/jwt-server/main.go

GOOSE_VERSION := latest
GOLANGCI_LINT_VERSION := v1.62.2
POSTGRES_DSN := "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable"

.PHONY: goose-up
goose-up:
	go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION) -dir ./migrations postgres $(POSTGRES_DSN) up

.PHONY: goose-down
goose-down:
	go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION) -dir ./migrations postgres $(POSTGRES_DSN) down

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

test:
	go test -v -count=1 ./... -coverprofile=cover.out

gen:
	go generate ./...

build:
	mkdir -p $(CURDIR)/bin && \
	go build -o $(CURDIR)/bin/jwt-server $(CURDIR)/cmd/jwt-server

docker-build:
	docker build -t jwt-server:local .

docker-compose-up: docker-build
	docker-compose up -d

swagger:
	docker run -p 8080:8080 -e SWAGGER_JSON=/openapi/openapi_v1.yml -v $(CURDIR):/openapi swaggerapi/swagger-ui