build:
	go build -o bin/tunhijack cmd/main.go

run: build
	sudo bin/tunhijack

test:
	sudo go test -v --v ./...

test-docker:
	sudo docker build -t tunhijack:test -f test.Dockerfile .