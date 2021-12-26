FROM golang:1.17.5-alpine as build
RUN apk --no-cache add ca-certificates
RUN mkdir /toby
WORKDIR /toby
COPY go.mod . 
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/toby
FROM scratch 
COPY --from=build /go/bin/toby /go/bin/toby
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/toby /toby
ENTRYPOINT ["/go/bin/toby"]

