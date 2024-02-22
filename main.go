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
	fmt.Println("Ingrese el comando a procesar: ")
	// Utilizando bufio.NewReader para leer la línea completa
	reader := bufio.NewReader(os.Stdin)
	command, _ = reader.ReadString('\n')
	//fmt.Scanln(&command)
	Analyzer.AnalyzeType(command)
}
