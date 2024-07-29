build-api:
	@go build -o cmd/bin/messager cmd/api/main.go

run-api:
	@./cmd/bin/messager

build-producer:
	@go build -o cmd/bin/producer cmd/producer/main.go

run-producer:
	@./cmd/bin/producer
