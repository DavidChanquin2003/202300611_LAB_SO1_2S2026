package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

// initValkey establece la conexión con la base de datos Valkey
func initValkey() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Prueba de conexión
	_, err := rdb.Ping(ctx).Result()
	return err
}

// saveLog guardará los eventos en Valkey para Grafana
func saveLog(mensaje string) {
	fmt.Println("-> Guardando en Valkey:", mensaje)
	// Aqui guardamos el log en una lista llamada "eventos_ebpf"
	rdb.LPush(ctx, "eventos_ebpf", mensaje)
}
