
.PHONY: test fmt

test:
	@go test -race -i ./...
	go test -race ./...

fmt:
	go fmt ./...