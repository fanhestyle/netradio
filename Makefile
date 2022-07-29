all: client server

client:
	go build ./cmd/client

server:
	go build ./cmd/server

.PHONY: clean

clean:
	$(RM) -rf server client