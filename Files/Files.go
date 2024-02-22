package Files

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

const PathFolder = "./MIA/P1/"

// Funtion to create bin file
func CreateFile(name string) error {
	//Ensure the directory exists
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err CreateFile dir==", err)
		return err
	}

	// Create file
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Err CreateFile create==", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

// Funtion to open bin file in read/write mode
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile==", err)
		return nil, err
	}
	return file, nil
}
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==", err)
		return err
	}
	return nil
}

// Function to Read an object from a bin file
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}

func ObtenerNuevoNombreArchivo() string {
	// Obtener la lista de archivos en el directorio
	archivos, err := ListArchivosCarpeta()
	if err != nil {
		return ""
	}
	const letras = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Calcular el siguiente nombre en orden alfabético
	nuevoNombre := letras[len(archivos) : len(archivos)+1]

	// Verificar que no se exceda la letra 'Z'
	if nuevoNombre > "Z" {
		fmt.Println("se excedió el límite de archivos")
		return ""
	}

	return nuevoNombre
}

func ListArchivosCarpeta() ([]string, error) {
	var archivos []string
	err := filepath.Walk(PathFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			archivos = append(archivos, info.Name())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return archivos, nil
}
