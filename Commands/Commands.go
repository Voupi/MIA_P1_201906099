package commands

import (
	Files "MIA_P1_201906099/Files"
	Models "MIA_P1_201906099/Models"
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
)

func MKDISK(nombreArchivo string, size *int, fit *string, unit *string) error { //mkdisk -size=3000 -unit=K  //mkdisk -size=10
	// Crear la estructura MBR
	mbr := Models.MBR{
		Tamano:        int32(*size),
		FechaCreacion: time.Now().Unix(),
		DiskSignature: int32(generateRandomSignature()),
	}
	copy(mbr.Fit[:], *fit)

	var fileName string
	// Crear el archivo binario
	nombreArchivo += ".dsk"
	fileName = Files.PathFolder + nombreArchivo
	Files.CreateFile(fileName)
	file, err := Files.OpenFile(fileName)
	if err != nil {
		return err
	}
	// Write 0 binary data to the file

	// create array of byte(0)
	for i := 0; i < *size; i++ {
		err := Files.WriteObject(file, byte(0), int64(i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
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
		var mbr Models.MBR
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

/* funcion que : RMDISK Elimina con rmdisk un archivo con el nombre que recibe del comando driveletter y que tenga mensaje
de confirmación de eliminación de disco*/

func RMDISK(fileName string) error { //rmdisk -driveletter=A
	// Confirmación de eliminación de disco
	// verificar si un archivo existe Recibe el path del archivo
	filePath := Files.PathFolder + fileName + ".dsk"
	fmt.Println("Es el path del archivo ", filePath)
	_, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("El archivo '%s' no existe.\n", fileName)
		return nil
	}
	if os.IsNotExist(err) { // Verificar si el archivo no existe
		fmt.Printf("El archivo '%s' no existe.\n", fileName)
		return err
	}
	fmt.Printf("¿Está seguro que desea eliminar el disco '%s'? (S/N): ", fileName) // Solicitar confirmación
	var confirmation string
	fmt.Scanln(&confirmation)
	if strings.ToUpper(confirmation) != "S" { // Verificar si la confirmación es "S"
		fmt.Println("Operación cancelada.")
		return nil
	}
	_, err = Files.OpenFile(filePath) // Abrir el archivo
	if err != nil {
		fmt.Printf("Error al leer el archivo '%s': %v\n", fileName, err)
		return err
	}
	err = os.Remove(filePath) // Eliminar el archivo
	if err != nil {
		fmt.Printf("Error al eliminar el archivo '%s': %v\n", fileName, err)
		return err
	}

	fmt.Printf("El disco '%s' ha sido eliminado exitosamente.\n", fileName)
	return nil
}

func FDISK(size int, driveletter string, name string, unit string, type_ string, fit string, delete string, add int) {
	fmt.Println("======Start FDISK======")
	fmt.Println("Size:", size)
	fmt.Println("Driveletter:", driveletter)
	fmt.Println("Name:", name)
	fmt.Println("Unit:", unit)
	fmt.Println("Type:", type_)
	fmt.Println("Fit:", fit)
	fmt.Println("Fit:")
	fmt.Println("Fit:", fit)
	filepath := Files.PathFolder + strings.ToUpper(driveletter) + ".dsk" // Obtener el path del archivo binario
	file, err := Files.OpenFile(filepath)                                // Abrir el archivo binario
	if err != nil {
		return
	}
	var TempMBR Models.MBR
	if err := Files.ReadObject(file, &TempMBR, 0); err != nil { // Lee el objeto del archivo binario
		return
	}
	Models.PrintMBR(TempMBR) // Imprime el objeto
	fmt.Println("-------------")

	var count = 0            // Contador de particiones
	var gap = int32(0)       // Espacio entre particiones
	var countPrimary = 0     // Contador de particiones primarias
	var countExtended = 0    // Contador de particiones extendidas
	for i := 0; i < 4; i++ { //Iterar sobre las particiones del disco
		if TempMBR.Partitions[i].Size != 0 { // Verificar si la partición tiene un tamaño o no
			count++
			gap = TempMBR.Partitions[i].Start + TempMBR.Partitions[i].Size // Calcular el espacio entre particiones
			if type_ == "p" {                                              //Validar que el tipo de partición sea primaria
				countPrimary++
			} else if type_ == "e" { //Validar que el tipo de partición sea extendida
				countExtended++
			}
		}
	} // Calcular el espacio entre particiones para la nueva partición
	if size > 0 { // Esto es para crear una nueva partición
		switch type_ {
		case "p":
			if countPrimary+countExtended >= 4 { // Verificar si ya existen 4 particiones primarias
				fmt.Printf("Error: Ya existen %d particiones primarias y %d extendidas\n", countPrimary, countExtended)
				return
			}
			for i := 0; i < 4; i++ {
				if TempMBR.Partitions[i].Size == 0 { //encontrar una partición vacía
					TempMBR.Partitions[i].Size = int32(size) // El tamaño de la partición es el tamaño ingresado

					if count == 0 { //Si no hay particiones
						TempMBR.Partitions[i].Start = int32(binary.Size(TempMBR)) //El inicio de la particion es el tamaño del MBR
					} else {
						TempMBR.Partitions[i].Start = gap // El inicio de la partición es el espacio entre particiones
					}

					copy(TempMBR.Partitions[i].Name[:], name)            // Copiar el nombre de la partición
					copy(TempMBR.Partitions[i].Fit[:], fit)              // Copiar el ajuste de la partición
					copy(TempMBR.Partitions[i].Status[:], "0")           // Copiar el estado de la partición
					copy(TempMBR.Partitions[i].Type[:], type_)           // Copiar el tipo de la partición
					TempMBR.Partitions[i].Correlative = int32(count + 1) // Copiar el número correlativo de la partición
					break
				}
			}
		case "e":
			if countExtended != 0 {
				fmt.Println("Error: Ya existe una partición extendida")
				return
			}
			for i := 0; i < 4; i++ {
				if TempMBR.Partitions[i].Size == 0 { //encontrar una partición vacía
					TempMBR.PartitionExtended[0].Size = int32(size) // El tamaño de la partición es el tamaño ingresado

					if count == 0 { //Si no hay particiones
						TempMBR.PartitionExtended[0].Start = int32(binary.Size(TempMBR)) //El inicio de la particion es el tamaño del MBR
					} else {
						TempMBR.PartitionExtended[0].Start = gap // El inicio de la partición es el espacio entre particiones
					}

					copy(TempMBR.PartitionExtended[0].Name[:], name)            // Copiar el nombre de la partición
					copy(TempMBR.PartitionExtended[0].Fit[:], fit)              // Copiar el ajuste de la partición
					copy(TempMBR.PartitionExtended[0].Status[:], "0")           // Copiar el estado de la partición
					copy(TempMBR.PartitionExtended[0].Type[:], type_)           // Copiar el tipo de la partición
					TempMBR.PartitionExtended[0].Correlative = int32(count + 1) // Copiar el número correlativo de la partición
					break
				}
			}
		case "l":
			if TempMBR.PartitionExtended[0].Size == 0 { // Verificar si ya existe una partición extendida
				fmt.Println("Error: No existe una partición extendida")
				return
			}
			var countLogical int         // Contador de particiones lógicas
			var gap = int32(0)           // Espacio entre particiones lógicas
			var espaciOcupado = int32(0) // Espacio ocupado por las particiones lógicas
			for j := 0; j < 20; j++ {    //Iterar sobre las particiones del disco
				if TempMBR.PartitionExtended[0].PartitionLogica[j].Size != 0 { // Verificar si la partición tiene un tamaño o no
					countLogical++
					gap = TempMBR.PartitionExtended[0].PartitionLogica[j].Start + TempMBR.PartitionExtended[0].PartitionLogica[j].Size // Calcular el espacio entre particiones
					espaciOcupado += TempMBR.PartitionExtended[0].PartitionLogica[j].Size
				}
			}
			if TempMBR.PartitionExtended[0].Size-espaciOcupado >= int32(size) { // Verificar si hay espacio suficiente
				fmt.Println("Error: No hay espacio suficiente para crear la partición lógica")
				return
			}
			for j := 0; j < 20; j++ { // Iterar sobre las particiones lógicas
				if TempMBR.PartitionExtended[0].PartitionLogica[j].Size == 0 { //encontrar una partición vacía
					TempMBR.PartitionExtended[0].PartitionLogica[j].Size = int32(size)
					TempMBR.PartitionExtended[0].PartitionLogica[j].Start = gap
					copy(TempMBR.PartitionExtended[0].PartitionLogica[j].Name[:], name)
					copy(TempMBR.PartitionExtended[0].PartitionLogica[j].Fit[:], fit)
					copy(TempMBR.PartitionExtended[0].PartitionLogica[j].Status[:], "0")
					TempMBR.PartitionExtended[0].PartitionLogica[j].Correlative = int32(countLogical + 1)
					break
				}
			}
			fmt.Println("Logical partition created successfully.")
		}

		if err := Files.WriteObject(file, TempMBR, 0); err != nil { // Sobreescribir el MBR
			return
		}
	}

	var TempMBR2 Models.MBR
	if err := Files.ReadObject(file, &TempMBR2, 0); err != nil { // Lee el objeto del archivo binario
		return
	}
	//Models.PrintMBR(TempMBR2)
	Models.PrintMBRWithExtended(TempMBR2)
	defer file.Close() // Cerrar el archivo binario
	fmt.Println("======End FDISK======")
}
