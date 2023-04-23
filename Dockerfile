FROM golang AS builder
WORKDIR /nlib
COPY go.mod /nlib/go.mod
COPY go.sum /nlib/go.sum
RUN go mod download
COPY . /nlib
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM node as builder-ui
WORKDIR /
RUN git clone https://github.com/borerer/nlib-dashboard.git
WORKDIR /nlib-dashboard
RUN npm install
RUN npm run build

FROM alpine
WORKDIR /nlib
COPY --from=builder /nlib/nlib /nlib/nlib
COPY --from=builder /nlib/config.yaml /nlib/config.yaml
COPY --from=builder-ui /nlib-dashboard/out /nlib/ui
ENTRYPOINT ["/nlib/nlib"]
EXPOSE 9502