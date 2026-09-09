package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func gestionarContenedores() {
	fmt.Println("Aplicando reglas de gestión de contenedores...")

	// Listamos solo los nombres de los contenedores activos con Docker
	out, err := exec.Command("docker", "ps", "--format", "{{.Names}}").Output()
	if err != nil {
		fmt.Println("Error al listar contenedores:", err)
		return
	}

	lineas := strings.Split(strings.TrimSpace(string(out)), "\n")
	var bajos, altos []string

	for _, nombre := range lineas {
		if nombre == "" || strings.Contains(nombre, "grafana") || strings.Contains(nombre, "valkey") {
			continue
		}

		// Clasificamos los contenedores por nombre
		if strings.HasPrefix(nombre, "intruso") {
			fmt.Println("Intruso detectado, eliminando de inmediato:", nombre)
			eliminarContenedor(nombre)
		} else if strings.HasPrefix(nombre, "bajo_consumo") {
			bajos = append(bajos, nombre)
		} else if strings.HasPrefix(nombre, "alto_") {
			altos = append(altos, nombre)
		}
	}

	// 3 contenedores de bajo consumo
	for len(bajos) > 3 {
		victima := bajos[0]
		fmt.Println("Limpiando exceso de bajo consumo:", victima)
		eliminarContenedor(victima)
		bajos = bajos[1:]
	}

	// 2 contenedores de alto consumo
	for len(altos) > 2 {
		victima := altos[0]
		fmt.Println("Limpiando exceso de alto consumo:", victima)
		eliminarContenedor(victima)
		altos = altos[1:]
	}
}

// eliminarContenedor fuerza el borrado. Esto envía un SIGKILL que atrapará eBPF
func eliminarContenedor(nombre string) {
	err := exec.Command("docker", "rm", "-f", nombre).Run()
	if err == nil {
		fmt.Println("   -> Contenedor eliminado:", nombre)
	}
}
