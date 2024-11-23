DEFAULT_PG_URL=postgres://user:password@localhost:5432/effective?sslmode=disable
COMPOSE=docker compose -f docker-compose.yml --env-file ./build/docker.env
COMPOSE_APP=${COMPOSE} -p app
COMPOSE_POSTGRES=${COMPOSE} -p postgres
COMPOSE_DEV=${COMPOSE} -p dev
LOCAL_BIN:=$(CURDIR)/bin

.PHONY: .up
up:
	${COMPOSE} up -d

.PHONY: .down
down:
	${COMPOSE} down

.PHONY: .up-postgres
up-postgres:
	${COMPOSE_POSTGRES} up -d --build postgres migrate-server

.PHONY: .down-postgres
down-postgres:
	${COMPOSE_POSTGRES} down

.PHONY: .up-dev
up-dev:
	${COMPOSE_DEV} up -d --build postgres migrate-server service

.PHONY: .down-dev
down-dev:
	${COMPOSE_DEV} down

.PHONY: .build-goose
build-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: .migration-up
migration-up:
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" up

.PHONY: .migration-down
migration-down:
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" down

.PHONY: .migration-status
migration-status:
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" status


.PHONY: .migration-create-sql
migration-create-sql: build-goose
	goose -dir ./migrations create $(filter-out $@,$(MAKECMDGOALS)) sql

.PHONY: unit-tests
unit-tests:
	go test  ./... -coverprofile coverage.txt

.PHONY: run-all-tests
run-all-tests: unit-tests
	$(eval PG_URL?=$(DEFAULT_PG_URL))
	TEST_DATABASE_URL=$(PG_URL)/effective go test ./tests/postgres/... -tags=integration

.PHONY: .generate-mockgen
generate-mockgen: generate-ifacemaker
	find . -name '*_mock.go' -delete
	go generate -x -run=mockgen ./internal/...

.PHONY: .bin-deps
bin-deps:
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install github.com/vburenin/ifacemaker@latest
	GOBIN=$(LOCAL_BIN) go install go.uber.org/mock/mockgen@latest

.PHONY: .generate-ifacemaker
generate-ifacemaker:
	$(LOCAL_BIN)/ifacemaker -f ./internal/usecase/song.go -s Song -i songUseCase -p mock_usecase -c "DONT EDIT: Auto generated" -o ./internal/usecase/mocks/song.go
	$(LOCAL_BIN)/ifacemaker -f ./internal/repo/song.go -s Song -i songRepo -p mock_repository -c "DONT EDIT: Auto generated" -o ./internal/repo/mocks/song.go
	$(LOCAL_BIN)/ifacemaker -f ./internal/repo/group.go -s Group -i groupStorage -p mock_repository -c "DONT EDIT: Auto generated" -o ./internal/repo/mocks/group.go
