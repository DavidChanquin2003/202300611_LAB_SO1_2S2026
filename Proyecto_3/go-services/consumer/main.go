package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type QueueMessage struct {
	Country         string `json:"country"`
	WarplanesInAir  int32  `json:"warplanes_in_air"`
	WarshipsInWater int32  `json:"warships_in_water"`
	Timestamp       string `json:"timestamp"`
}

func main() {
	// 1. Conectar a Valkey (usamos el nombre del Service creado por tu Operador)
	rdb := redis.NewClient(&redis.Options{
		Addr: "valkeyinstance-sample-valkey-svc:6379",
		DB:   0,
	})
	ctx := context.Background()

	// Probar conexión a Valkey
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Error conectando a Valkey: %v", err)
	}
	log.Println("Conexión exitosa a la base de datos Valkey")

	// 2. Conectar a RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq-svc:5672/")
	if err != nil {
		log.Fatalf("Error conectando a RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error abriendo canal RabbitMQ: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("war_reports", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error declarando cola: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error registrando consumidor: %v", err)
	}

	forever := make(chan os.Signal, 1)
	signal.Notify(forever, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for d := range msgs {
			var msg QueueMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("Error decodificando mensaje: %v", err)
				continue
			}

			log.Printf("Reporte leído de RabbitMQ: País=%s", msg.Country)

			// 3. Guardar el reporte en Valkey como un Stream con campos nativos
			err := rdb.XAdd(ctx, &redis.XAddArgs{
				Stream: "war_reports_stream",
				Values: map[string]interface{}{
					"country":           msg.Country,
					"warplanes_in_air":  msg.WarplanesInAir,
					"warships_in_water": msg.WarshipsInWater,
					"timestamp":         msg.Timestamp,
				},
			}).Err()

			if err != nil {
				log.Printf("Error guardando en Valkey: %v", err)
			} else {
				log.Printf("Reporte guardado en Valkey exitosamente")
			}
		}
	}()

	log.Println("Consumidor Go iniciado. Esperando reportes...")
	<-forever
}
