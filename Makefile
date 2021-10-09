build:
	go build -o target/escpserver.
lint:
	golangci-lint run . internal/...
clean:
	rm -rf target/