package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// la chamba de analizarProcesos es que lee el archivo virtual creado por el módulo en C
func analizarProcesos() {
	filePath := "/proc/continfo_pr2_so1_202300611"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error leyendo %s. ¿Está cargado el módulo del kernel?\n", filePath)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fmt.Println("Extrayendo métricas desde el Kernel...")

	// se lee línea por línea
	for scanner.Scan() {
		line := scanner.Text()

		// aqui solo consultamos las metricas de la RAM, las que nos interesan
		if strings.Contains(line, "Total de memoria RAM") ||
			strings.Contains(line, "Memoria RAM libre") ||
			strings.Contains(line, "Memoria RAM en uso") {

			// se envia esta estadística a Valkey
			saveLog(line)
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error al escanear el archivo: %v\n", err)
	}
}
