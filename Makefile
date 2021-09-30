build:
	CGO_ENABLED=0 go build -o bin/tunhijack cmd/main.go

run: build
	sudo bin/tunhijack

test:
	sudo go test -v --v ./...

test-docker:
	sudo docker build -t tunhijack:test -f test.Dockerfile .
	sudo docker run --device /dev/net/tun --cap-add NET_ADMIN --cap-add SYS_ADMIN tunhijack:test