package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Iniciando consumo de RAM...")
	var memoryBlock [][]byte
	for {
		// Le puse 10MB de memoria en cada iteración
		block := make([]byte, 10*1024*1024)
		for i := range block {
			block[i] = 1
		}
		memoryBlock = append(memoryBlock, block)
		time.Sleep(2 * time.Second)
	}
}
