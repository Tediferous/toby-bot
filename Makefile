default:
	go build -o bin/toby

install:
	go build -o "${GOPATH}/bin/toby"

fmt:
	go fmt ./...
