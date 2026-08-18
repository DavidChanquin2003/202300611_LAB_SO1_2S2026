# Construcción
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o api3 main.go

# Imagen final ligera
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api3 .
EXPOSE 8083
CMD ["./api3"]