PWD=$(shell pwd)
GO_GENERATE_DEPS=$(shell grep --recursive --files-with-matches 'go:generate' .)
SQLC_DEPS=$(shell ls db/sqlc/*.sql | grep --invert-match db/sqlc/schema.sql)
DEFAULT_GO_VERSION=$(shell cat go.mod | grep 'go \d.\d\d' | sed 's/go //')
NON_INTEGRATION_TESTS=$(shell go list ./... | grep -v test/)
INTEGRATION_TESTS=$(shell go list ./... | grep test/integration)

GO_VERSION?=${DEFAULT_GO_VERSION}
POSTGRES_VERSION?=13.2
SQLC_VERSION?=1.10.0
PGDATABASE?=poc
PGHOST?=pocdb
PGPASSWORD?=changeme
PGPORT?=5432
PGUSER?=poc
PGADMIN_DEFAULT_EMAIL?=admin@admin.com
PGADMIN_DEFAULT_PASSWORD?=admin123
PGUSER?=poc
DB_CONN_STRING?=postgres://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}
HTTP_LISTEN_ADDRESS?=0.0.0.0:8321

.PHONY: test-integration
test-integration: start-services build/.empty-targets/generate
	docker run \
		--interactive \
		--tty \
		--network poc-demo \
		--env DB_CONN_STRING=${DB_CONN_STRING} \
		--env HTTP_LISTEN_ADDRESS=${HTTP_LISTEN_ADDRESS} \
		--volume ${PWD}:/go/src/github.com/slcjordan/poc \
		--workdir /go/src/github.com/slcjordan/poc \
		golang:${GO_VERSION} sh -c 'go test -v ${INTEGRATION_TESTS}'

.PHONY: test
test: build/.empty-targets/generate
	docker run \
		--interactive \
		--tty \
		--network poc-demo \
		--publish 8411:8411 \
		--volume ${PWD}:/go/src/github.com/slcjordan/poc \
		--workdir /go/src/github.com/slcjordan/poc \
		golang:${GO_VERSION} sh -c 'go test -v ${NON_INTEGRATION_TESTS}'

.PHONY: admin
admin:
	docker run \
		--interactive \
		--tty \
		--env PGADMIN_DEFAULT_EMAIL=${PGADMIN_DEFAULT_EMAIL} \
		--env PGADMIN_DEFAULT_PASSWORD=${PGADMIN_DEFAULT_PASSWORD} \
		--network poc-demo \
		--publish 8180:80 \
		--volume ${PWD}/data/pgadmin:/var/lib/pgadmin \
		dpage/pgadmin4

.PHONY: run-docs
run-docs:
	echo 'once running, please visit http://localhost:8411/pkg/github.com/slcjordan/poc/ for documentation'
	docker run \
		--interactive \
		--tty \
		--network poc-demo \
		--publish 8411:8411 \
		--volume ${PWD}:/go/src/github.com/slcjordan/poc \
		--workdir /go/src/github.com/slcjordan/poc/cmd/api \
		golang:${GO_VERSION} sh -c 'go get -v golang.org/x/tools/cmd/godoc && godoc -http=:8411'

.PHONY: generate
generate: build/.empty-targets/generate

.PHONY: sqlc
sqlc: build/.empty-targets/sqlc

.PHONY: stop-services
stop-services:
	docker stop pocdb

.PHONY: start-services
start-services: build/.empty-targets/network
	docker container inspect --format='database is {{.State.Status}}' pocdb || docker run \
		--detach \
		--name pocdb \
		--rm \
		--env POSTGRES_PASSWORD=${PGPASSWORD} \
		--env POSTGRES_USER=${PGUSER} \
		--env PGDATA=/var/lib/postgresql/data/pgdata \
		--network 'poc-demo' \
		--volume ${PWD}/data:/var/lib/postgresql/data \
		postgres:${POSTGRES_VERSION}

build/.empty-targets/generate: ${GO_GENERATE_DEPS}
	@mkdir -p test
	@echo "(re)generating mocks"
	- rm test/mocks/*
	docker run \
		--interactive \
		--rm \
		--volume $(PWD):/go/src/github.com/slcjordan/poc \
		--workdir /go/src/github.com/slcjordan/poc \
		golang:1.16 \
		sh -c 'go get -v golang.org/x/tools/cmd/stringer && go get -v github.com/golang/mock/mockgen@v1.6.0 && go generate -v ./...'
	@mkdir -p $(@D)
	@touch $@

build/.empty-targets/network:
	docker network create poc-demo
	@mkdir -p $(@D)
	@touch $@

.PHONY: psql
psql: start-services
	docker run \
		--interactive \
		--tty \
		--publish 5432:5432 \
		--network poc-demo \
		--env PGDATABASE=${PGDATABASE} \
		--env PGHOST=${PGHOST} \
		--env PGPASSWORD=${PGPASSWORD} \
		--env PGPORT=${PGPORT} \
		--env PGUSER=${PGUSER} \
		--workdir / \
		postgres:${POSTGRES_VERSION} psql

.PHONY: run-dev
run-dev: start-services
	docker run \
		--interactive \
		--tty \
		--network poc-demo \
		--publish 8321:8321 \
		--env DB_CONN_STRING=${DB_CONN_STRING} \
		--env HTTP_LISTEN_ADDRESS=${HTTP_LISTEN_ADDRESS} \
		--volume ${PWD}:/go/src/github.com/slcjordan/poc \
		--workdir /go/src/github.com/slcjordan/poc/cmd/api \
		golang:${GO_VERSION} go run main.go

db/sqlc/schema.sql: start-services
	docker run \
		--interactive \
		--tty \
		--network poc-demo \
		--env PGDATABASE=${PGDATABASE} \
		--env PGHOST=${PGHOST} \
		--env PGPASSWORD=${PGPASSWORD} \
		--env PGPORT=${PGPORT} \
		--env PGUSER=${PGUSER} \
		--volume ${PWD}/db:/db \
		--workdir / \
		postgres:${POSTGRES_VERSION} pg_dump \
			--file $@ \
			--schema-only

build/.empty-targets/sqlc: ${SQLC_DEPS}
	docker run \
		--interactive \
		--tty \
		--network poc-demo \
		--volume ${PWD}:/repo \
		--workdir /repo \
		kjconroy/sqlc:${SQLC_VERSION} generate
	@mkdir -p $(@D)
	@touch $@
