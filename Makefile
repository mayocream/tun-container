build:
	CGO_ENABLED=0 go build -o bin/tunhijack cmd/main.go

run: build
	sudo bin/tunhijack

test:
	sudo go test -v --v ./...

test-docker:
	sudo docker build -t tunhijack:test -f test.Dockerfile .
	sudo docker run --device /dev/net/tun --cap-add NET_ADMIN --cap-add SYS_ADMIN tunhijack:test

tun2socks:
	tun2socks -device tun0 -proxy socks5://127.0.0.1:1234

ipconfig:
	ip addr add 10.0.0.1/24 dev tun0
	ip link set tun0 up