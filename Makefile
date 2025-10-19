build:
	pnpm build && go build -o bin/sl-sime-table ./cmd/webserver/main.go
run:
	IS_DEV=true go run ./cmd/webserver/main.go
test:
	go test ./...
test-race:
	go test ./... -race
test-verbose:
	go test ./... -v
clean:
	rm -rf bin/
