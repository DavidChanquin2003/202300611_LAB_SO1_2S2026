# Etapa 1: Compilación
FROM rust:slim-bullseye AS builder
WORKDIR /app
# Copiamos los archivos de dependencias y el código fuente
COPY Cargo.toml Cargo.lock ./
COPY src ./src
# Compilamos en modo release para máxima velocidad
RUN cargo build --release

# Etapa 2: Imagen final ligera
FROM debian:bullseye-slim
WORKDIR /app
# Copiamos solo el binario compilado de la etapa anterior
COPY --from=builder /app/target/release/api-rust /usr/local/bin/api-rust
# Exponemos el puerto
EXPOSE 8080
# Comando de inicio
CMD ["api-rust"]