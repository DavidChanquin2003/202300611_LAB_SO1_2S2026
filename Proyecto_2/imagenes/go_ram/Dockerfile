FROM golang:1.20-alpine
WORKDIR /app
COPY main.go .
RUN go build -o ram_consumer main.go
CMD ["./ram_consumer"]