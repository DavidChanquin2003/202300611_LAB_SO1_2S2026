package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Iniciando Daemon Gestor de Contenedores SO1...")

	// Conectar a Valkey
	if err := initValkey(); err != nil {
		fmt.Printf("Error conectando a Valkey: %v\n", err)
		return
	}
	fmt.Println("Conectado a Valkey exitosamente.")

	// Levantar el Cronjob
	fmt.Println("Configurando Cronjob para generador.sh...")
	setupCronjob()

	fmt.Println("Ejecutando script para cargar módulo del kernel...")
	cmd := exec.Command("sudo", "bash", "../cargar_modulo.sh")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Nota: El módulo ya estaba cargado o hubo un detalle (%v)\n", err)
	}

	// Manejo de cierre seguro (o sea cuando se presiona Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup() // Se limpia el cronjob antes de salir
		os.Exit(0)
	}()

	go iniciarEscuchaEBPF()

	// Loop Principal
	for {
		fmt.Println("\n--- Ejecutando ciclo de análisis (Esperate 20s) ---")

		analizarProcesos()

		fmt.Println("Esperando 20 segundos para el siguiente ciclo...")
		time.Sleep(20 * time.Second) // Ejecución cada 20 segundos
	}
}

// setupCronjob agrega el script generador al crontab del sistema
func setupCronjob() {
	cwd, _ := os.Getwd()
	scriptPath := cwd + "/../generador.sh"

	// El comando agrega la tarea al crontab para ejecutarse cada 2 minutos
	cmdStr := fmt.Sprintf("(crontab -l 2>/dev/null; echo \"*/2 * * * * %s\") | crontab -", scriptPath)
	exec.Command("bash", "-c", cmdStr).Run()
}

// cleanup elimina el cronjob al apagar el servicio
func cleanup() {
	fmt.Println("\nApagando Daemon de Go de forma segura...")
	// Se borra el crontab del usuario actual
	exec.Command("bash", "-c", "crontab -r").Run()
	fmt.Println("Cronjob eliminado. ¡Adiós!")
}
