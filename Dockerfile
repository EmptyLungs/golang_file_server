FROM golang:1.21-alpine3.18 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin /app/cmd/api/*

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/bin /app/bin
USER nobody
EXPOSE 8080
CMD ["/app/bin"]
