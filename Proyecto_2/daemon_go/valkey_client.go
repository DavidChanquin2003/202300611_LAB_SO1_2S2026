package main

import (
	"context"
	"fmt"
	"strings"

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

	if strings.Contains(mensaje, "Total de memoria RAM") {
		rdb.LPush(ctx, "ram_total", mensaje)
	} else if strings.Contains(mensaje, "Memoria RAM libre") {
		rdb.LPush(ctx, "ram_libre", mensaje)
	} else if strings.Contains(mensaje, "Memoria RAM en uso") {
		rdb.LPush(ctx, "ram_en_uso", mensaje)
	} else if strings.Contains(mensaje, "eBPF") {
		// Si es evento de eBPF, a su lista propia
		rdb.LPush(ctx, "eventos_ebpf", mensaje)
	}

	rdb.LPush(ctx, "log_general", mensaje)
}
