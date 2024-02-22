package commands

import (
	Files "MIA_P1_201906099/Files"
	MBR "MIA_P1_201906099/Models"
	"bufio"
	"fmt"
	"strings"
	"time"
)

func GenerarDiscoBinario(nombreArchivo string) error {
	// Crear la estructura MBR
	mbr := MBR.MBR{
		Tamano:        5 * 1024 * 1024,
		FechaCreacion: time.Now().Unix(),
		DiskSignature: int32(generateRandomSignature()),
	}
	var fileName string
	// Crear el archivo binario
	nombreArchivo += ".sdk"
	fileName = Files.PathFolder + nombreArchivo
	Files.CreateFile(fileName)
	file, err := Files.OpenFile(fileName)
	if err != nil {
		return err
	}
	Files.WriteObject(file, mbr, 0)
	defer file.Close()
	fmt.Printf("Archivo binario '%s' creado exitosamente.\n", nombreArchivo)
	return nil
}

func generateRandomSignature() int {
	return int(time.Now().UnixNano())
}
func Rep() error {
	// Revisar archivos en la carpeta "MIA/P1/"
	files, err := Files.ListArchivosCarpeta()
	if err != nil {
		return err
	}
	for _, fileName := range files {
		fmt.Println(fileName)

		// Abrir el archivo binario
		// Open bin file
		file, err := Files.OpenFile(Files.PathFolder + fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		var mbr MBR.MBR
		if err := Files.ReadObject(file, &mbr, 0); err != nil {
			return err
		}
		file.Close()
		// Imprimir datos del MBR
		fmt.Println("El nombre del archivo es: ", fileName)
		fmt.Println("Datos del MBR:")
		fmt.Printf("Tamaño total del disco: %d bytes\n", mbr.Tamano)
		fmt.Printf("Fecha y hora de creación: %s\n", time.Unix(mbr.FechaCreacion, 0))
		fmt.Printf("Disk Signature: %d\n", mbr.DiskSignature)
	}
	return nil
}
func Execute(path string, lineas *[]string) error {
	scannerPath := bufio.NewScanner(strings.NewReader(path))
	scannerPath.Scan()
	path = scannerPath.Text()
	file, err := Files.OpenFile(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// Crear un escáner para leer el archivo línea por línea
	scanner := bufio.NewScanner(file)

	// Recorrer el archivo línea por línea
	for scanner.Scan() {
		linea := scanner.Text()
		if linea[0:1] != "#" {
			fmt.Println("No es Comentario ")
			*lineas = append(*lineas, linea)
		}
	}
	// Verificar si hubo errores durante la lectura del archivo
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
