# Building stage for Marty application
FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app cmd/gophermart/main.go

# Final stage
FROM alpine:latest
COPY --from=builder /app/app /app/
CMD ["/app/app", "-a", "localhost:8081", "-d", "postgres://user:password@dbhost:5432/dbname", "-r", "http://localhost:8080"]
