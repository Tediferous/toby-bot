default:
	go build main.go

install:
	go build -o "${GOPATH}/bin/toby"

fmt:
	go fmt ./...
