build:
	@go build -o ./bin/server.go
run: build
	@./bin/server.go
