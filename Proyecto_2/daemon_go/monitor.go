package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Estructura para almacenar los datos de cada proceso
type Proceso struct {
	PID  string
	Name string
	VSZ  string
	RSS  string
	RAM  string
	CPU  string
}

// analizarProcesos lee el archivo virtual creado por el módulo en C
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

	var procesos []Proceso

	for scanner.Scan() {
		line := scanner.Text()

		// Se guarda la info de la RAM en Valkey
		if strings.Contains(line, "Total de memoria RAM") || strings.Contains(line, "Memoria RAM libre") || strings.Contains(line, "Memoria RAM en uso") {
			saveLog(line)
			continue
		}

		// Deserialización de los procesos
		if strings.Contains(line, "|") && !strings.Contains(line, "PID") {
			parts := strings.Split(line, "|")
			if len(parts) >= 6 {
				p := Proceso{
					PID:  strings.TrimSpace(parts[0]),
					Name: strings.TrimSpace(parts[1]),
					VSZ:  strings.TrimSpace(parts[2]),
					RSS:  strings.TrimSpace(parts[3]),
					RAM:  strings.TrimSpace(parts[4]),
					CPU:  strings.TrimSpace(parts[5]),
				}
				procesos = append(procesos, p)

				// Filtro
				if p.Name == "ram_consumer" || p.Name == "sleep" || p.Name == "sh" || p.Name == "cat" {
					fmt.Printf("Detectado -> Nombre: %-12s | PID: %-6s | RAM: %-4s | CPU (Jiffies): %s\n", p.Name, p.PID, p.RAM, p.CPU)
				}
			}
		}
	}

	fmt.Printf("Se analizaron un total de %d procesos del sistema.\n", len(procesos))

	// TOP 5 RAM
	sort.Slice(procesos, func(i, j int) bool {
		ramI, _ := strconv.ParseFloat(strings.ReplaceAll(procesos[i].RAM, "%", ""), 64)
		ramJ, _ := strconv.ParseFloat(strings.ReplaceAll(procesos[j].RAM, "%", ""), 64)
		return ramI > ramJ // Mayor a menor
	})

	clearRanking("top_ram")
	for i := 0; i < len(procesos) && i < 5; i++ {
		valorRam, _ := strconv.ParseFloat(strings.ReplaceAll(procesos[i].RAM, "%", ""), 64)
		saveRanking("top_ram", valorRam, procesos[i].Name+" ("+procesos[i].PID+")")
	}

	sort.Slice(procesos, func(i, j int) bool {
		cpuI, _ := strconv.ParseFloat(procesos[i].CPU, 64)
		cpuJ, _ := strconv.ParseFloat(procesos[j].CPU, 64)
		return cpuI > cpuJ
	})

	clearRanking("top_cpu")
	for i := 0; i < len(procesos) && i < 5; i++ {
		valorCpu, _ := strconv.ParseFloat(procesos[i].CPU, 64)
		saveRanking("top_cpu", valorCpu, procesos[i].Name+" ("+procesos[i].PID+")")
	}
	gestionarContenedores()
}
