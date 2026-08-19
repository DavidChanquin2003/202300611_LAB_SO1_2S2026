package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const CARNET = "202300611"

type HealthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	VM        int    `json:"VM"`
	Carnet    string `json:"carnet"`
}

type CallResponse struct {
	Apiname    string `json:"apiname"`
	Message    string `json:"message"`
	Connection bool   `json:"connection"`
	Carnet     string `json:"carnet"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := HealthResponse{
		Status:    "UP",
		Message:   "API2 is Ready",
		Timestamp: time.Now().Format(time.RFC3339),
		VM:        1, // la API2 también va en la VM1
		Carnet:    CARNET,
	}
	json.NewEncoder(w).Encode(response)
}

func callApiHandler(targetApiName string, targetUrl string, targetVM int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(targetUrl)

		success := false
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
	http.HandleFunc("/health", healthHandler)

	// Este es el endpoint para llamar a la API1
	api1Url := "http://192.168.122.196:8081/health"
	http.HandleFunc("/api2/"+CARNET+"/call-api1", callApiHandler("API1", api1Url, 1))

	// Este es el endpoint para llamar a la API3
	api3Url := "http://192.168.122.77:8083/health"
	http.HandleFunc("/api2/"+CARNET+"/call-api3", callApiHandler("API3", api3Url, 2))

	fmt.Println("API2 inicializada correctamente. Escuchando en el puerto 8082...")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
