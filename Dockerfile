FROM golang:1.25-alpine3.21 AS modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:1.25-alpine3.21 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0  \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o /bin/app ./cmd/app

FROM alpine:3.21
COPY --from=builder /app/configs /configs
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
CMD ["/app"]