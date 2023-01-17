FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o /usr/bin/example-service ./cmd/example-service/main.go

FROM builder as tester
RUN go test ./...

FROM alpine as runner
COPY --chown=0:0 --from=builder /usr/bin/example-service /usr/bin
ENTRYPOINT exec example-service
