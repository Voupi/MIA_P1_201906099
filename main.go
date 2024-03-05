package main

import (
	"MIA_P1_201906099/Analyzer"
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("===Start=== Usuario: Daniel Chan   Carné: 201906099")
	var command string
	for {
		fmt.Println("Ingrese el comando a procesar (o 'salir' para terminar): ")
		reader := bufio.NewReader(os.Stdin)
		command, _ = reader.ReadString('\n')
		if command == "salir\n" {
			break
		}
		Analyzer.AnalyzeType(command)
	}
}
