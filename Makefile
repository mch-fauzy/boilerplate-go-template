BINARY=engine
test: clean documents generate
	go test -v -cover -covermode=atomic ./...

coverage: clean documents generate
	bash coverage.sh --html

dev: generate
	go run github.com/cosmtrek/air

run: generate
	go run .

build:
	go build -o ${BINARY} .

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	@find . -name *mock* -delete
	@rm -rf .cover wire_gen.go docs

docker_build:
	docker build -t boilerplate-go -f Dockerfile-local .

docker_start:
	docker-compose up --build

docker_stop:
	docker-compose down

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

generate:
	go generate ./...
	
.PHONY: test coverage engine clean build docker run stop lint-prepare lint documents generate
