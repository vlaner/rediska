build:
	@go build -o bin/main cmd/main.go

run:
	@go run cmd/main.go

test:
	go test -cover -v ./...

cover:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html 
	rm cover.out