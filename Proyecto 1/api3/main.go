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
		Message:   "API3 is Ready",
		Timestamp: time.Now().Format(time.RFC3339),
		VM:        2, // la API3 va en la VM2 (la que es para Podman)
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

	// Llamada a la API1 (VM1)
	api1Url := "http://192.168.122.196:8081/health"
	http.HandleFunc("/api3/"+CARNET+"/call-api1", callApiHandler("API1", api1Url, 1))

	// Llamada a la API2 (VM1)
	api2Url := "http://192.168.122.196:8082/health"
	http.HandleFunc("/api3/"+CARNET+"/call-api2", callApiHandler("API2", api2Url, 1))

	fmt.Println("API3 inicializada correctamente. Escuchando en el puerto 8083...")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
