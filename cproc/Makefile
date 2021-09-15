.PHONY: examples
examples: ping sql

.PHONY: ping
ping:
	go build -o build/ ./examples/ping

.PHONY: sql
	go build -o build/ ./examples/sql

.PHONY: test
test:
	go test -cover ./...