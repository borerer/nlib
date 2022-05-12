FROM golang:bullseye AS builder
WORKDIR /nlib
COPY go.mod /nlib/go.mod
COPY go.sum /nlib/go.sum
RUN go mod download
COPY . /nlib
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM alpine
WORKDIR /nlib
COPY --from=builder /nlib/nlib /nlib/nlib
COPY --from=builder /nlib/data/config.json /nlib/data/config.json
ENTRYPOINT ["/nlib/nlib"]
EXPOSE 9502