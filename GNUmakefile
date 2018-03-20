TEST?=$$(go list ./... |grep -Ev 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: check test

build:: ## Build the API server
	mkdir -p ./bin
	govvv build -o bin/triton-sg ./cmd/triton-sg

tools:: ## Download and install all dev/code tools
	@echo "==> Installing dev/build tools"
	go get -u github.com/ahmetb/govvv
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

test:: ## Run unit tests
	@echo "==> Running unit test with coverage"
	@./scripts/go-test-with-coverage.sh

testacc:: ## Run acceptance tests
	@echo "==> Running acceptance tests"
	TRITON_TEST=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

check::
	gometalinter \
			--deadline 10m \
			--vendor \
			--sort="path" \
			--aggregate \
			--enable-gc \
			--disable-all \
			--enable goimports \
			--enable misspell \
			--enable vet \
			--enable deadcode \
			--enable varcheck \
			--enable ineffassign \
			--enable errcheck \
			--enable gofmt \
			./...

dev-db-start:: ## Start the development database
	@echo "==> Running docker-compose up"
	docker-compose up -d

dev-db-stop:: ## Stop the development database
	@echo "==> Running docker-compose kill"
	docker-compose kill
	rm -rf data/

dev-db-clean:: ## Cleans CRDB of all data
	cockroach sql --database triton --host localhost --insecure --certs-dir ./dev/vagrant/certs < ./dev/clean.sql

dev-db-seed:: ## Seed CRDB with test data (run from project root)
	cockroach sql --database triton --host localhost --insecure --certs-dir ./dev/vagrant/certs < ./dev/backup.sql

.PHONY: help
help:: ## Display this help message
	@echo "GNU make(1) targets:"
	@grep -E '^[a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
