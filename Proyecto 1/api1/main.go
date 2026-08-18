package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const CARNET = "202300611"

// Estructura para el end point /health
type HealthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	VM        int    `json:"VM"`
	Carnet    string `json:"carnet"`
}

// forma de respuestas para las llamadas a otras APIs
type CallResponse struct {
	Apiname    string `json:"apiname"`
	Message    string `json:"message"`
	Connection bool   `json:"connection"`
	Carnet     string `json:"carnet"`
}

// Controlador para el endpoint /health
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := HealthResponse{
		Status:    "UP",
		Message:   "API1 is Ready",
		Timestamp: time.Now().Format(time.RFC3339),
		VM:        1,
		Carnet:    CARNET,
	}

	json.NewEncoder(w).Encode(response)
}

// este es el controlador para llamar a otras APIs y verificar su estado
func callApiHandler(targetApiName string, targetUrl string, targetVM int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(targetUrl)

		success := false

		// Verificamos si la petición fue exitosa y si el status es "UP"
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				var healthResp HealthResponse
				bodyBytes, _ := io.ReadAll(resp.Body)
				json.Unmarshal(bodyBytes, &healthResp)

				if healthResp.Status == "UP" {
					success = true
				}
			}
		}

		// respuesta final si la respuesta fue exitosa o nel
		var callResp CallResponse
		callResp.Apiname = targetApiName
		callResp.Carnet = CARNET

		if success {
			callResp.Message = fmt.Sprintf("The %s located on the VM%d is working", targetApiName, targetVM)
			callResp.Connection = true
		} else {
			callResp.Message = fmt.Sprintf("ERROR: The %s located on the VM%d is not working", targetApiName, targetVM)
			callResp.Connection = false
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(callResp)
	}
}

func main() {
	// Endpoint principal
	http.HandleFunc("/health", healthHandler)

	// Endpoint de llamada 1 (hacia API2)
	api2Url := "http://192.168.122.196:8082/health"
	http.HandleFunc("/api1/"+CARNET+"/call-api2", callApiHandler("API2", api2Url, 1))

	// Endpoint de llamada 2 (hacia API3 en la VM2)
	api3Url := "http://192.168.122.77:8083/health"
	http.HandleFunc("/api1/"+CARNET+"/call-api3", callApiHandler("API3", api3Url, 2))

	// Inicio el servidor
	fmt.Println("API1 inicializada correctamente. Escuchando en el puerto 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
