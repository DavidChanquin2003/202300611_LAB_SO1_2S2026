package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// iniciarEscuchaEBPF corre el script de Python y lee sus resultados en tiempo real
func iniciarEscuchaEBPF() {
	fmt.Println("Iniciando conexión con sonda eBPF en segundo plano...")

	cmd := exec.Command("sudo", "python3", "../ebpf/monitor_kill.py")

	// Esto es un camino para leer lo que Python imprime en la terminal
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error al conectar con eBPF:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error al arrancar la sonda eBPF:", err)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		linea := scanner.Text()

		if strings.Contains(linea, "9") || strings.Contains(linea, "15") {
			fmt.Println(" [eBPF Confirmación Kernel]:", linea)

			saveLog("Contenedor Eliminado | Datos eBPF: " + linea)
		}
	}
}
