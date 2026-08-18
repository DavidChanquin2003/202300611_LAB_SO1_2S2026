# Construcción
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o api2 main.go

# Imagen final ligera
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api2 .
EXPOSE 8082
CMD ["./api2"]