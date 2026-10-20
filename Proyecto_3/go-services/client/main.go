package main

import (
	"context"
	"log"
	"time"

	pb "client/proto" // Importamos los protos locales

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Estructura que recibiremos desde Rust
type WarReport struct {
	Country         string `json:"country"`
	WarplanesInAir  int32  `json:"warplanes_in_air"`
	WarshipsInWater int32  `json:"warships_in_water"`
	Timestamp       string `json:"timestamp"`
}

var countryMap = map[string]pb.Countries{
	"USA": pb.Countries_usa,
	"RUS": pb.Countries_rus,
	"CHN": pb.Countries_chn,
	"ESP": pb.Countries_esp,
	"GTM": pb.Countries_gtm,
}

func main() {
	app := fiber.New()

	// Conexión gRPC hacia el Deployment 2 (en K8s lo llamaremos "go-server-svc")
	conn, err := grpc.Dial("go-server-svc:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar a gRPC: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewWarReportServiceClient(conn)

	app.Post("/send", func(c *fiber.Ctx) error {
		var req WarReport
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "JSON inválido"})
		}

		// Mapear el string del país al Enum de Protobuf
		countryEnum := pb.Countries_countries_unknown
		if val, ok := countryMap[req.Country]; ok {
			countryEnum = val
		}

		grpcReq := &pb.WarReportRequest{
			Country:         countryEnum,
			WarplanesInAir:  req.WarplanesInAir,
			WarshipsInWater: req.WarshipsInWater,
			Timestamp:       req.Timestamp,
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		// Enviar al Servidor Go mediante gRPC
		res, err := grpcClient.SendReport(ctx, grpcReq)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error enviando a gRPC"})
		}

		return c.JSON(fiber.Map{
			"status":  res.Status,
			"message": res.Message,
		})
	})

	log.Println("Cliente Go escuchando peticiones HTTP en puerto 3000...")
	app.Listen(":3000")
}
