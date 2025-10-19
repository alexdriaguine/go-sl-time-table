build:
	go build -o bin/sl-sime-table ./cmd/webserver/main.go
run:
	go run ./cmd/webserver/main.go
test:
	go test ./...
test-race:
	go test ./... -race
test-verbose:
	go test ./... -v
clean:
	rm -rf bin/
