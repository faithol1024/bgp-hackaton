.PHONY: build-http
build-http:
	go build -v -o bin/app-http cmd/app-http/*.go

.PHONY: run-http
run-http:
	go build -v -o bin/app-http cmd/app-http/*.go
	@./bin/app-http

.PHONY: test
test:
	@which gotest 2>/dev/null || go get -v github.com/rakyll/gotest
	@gotest -v --race ./...

.PHONY: lint
lint:
	@which golangci-lint 2>/dev/null || go get -v -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@golangci-lint run ./... --disable errcheck

.PHONY: docker-up
docker-up: 
	docker-compose up --build -d

.PHONY: docker-down
docker-down:
	docker-compose down