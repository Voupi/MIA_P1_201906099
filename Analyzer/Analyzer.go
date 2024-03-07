package Analyzer

import (
	Commands "MIA_P1_201906099/Commands"
	"MIA_P1_201906099/Files"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

func AnalyzeType(command string) { //Método donde se va a redirificar el comando a su respectiva función
	if strings.Contains(command, "mkdisk") {
		processMkdiskCommand(command)
	} else if strings.Contains(command, "EXECUTE") {
		processExecuteCommand(command)
	} else if strings.Contains(command, "REP") || strings.Contains(command, "rep") {
		fmt.Println("La cadena contiene la palabra 'REP'.")
		Commands.Rep()
	} else if strings.Contains(command, "RMDISK") || strings.Contains(command, "rmdisk") {
		// Define flags
		fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
		driveletter := fs.String("driveletter", "", "Letra")
		fs.Parse(os.Args[1:])
		matches := re.FindAllStringSubmatch(command, -1)
		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "driveletter", "name":
				fs.Set(flagName, flagValue)
			default:
				fmt.Println("Error: Flag not found")
			}
		}
		Commands.RMDISK(*driveletter)
	} else if strings.Contains(command, "FDISK") || strings.Contains(command, "fdisk") {
		processFDISKCommand(command)
	}
}
func processFDISKCommand(command string) {
	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño") //definiendo el tamaño como 0 predeterminadamente
	driveletter := fs.String("driveletter", "", "Letra") //definiendo la letra como vacío predeterminadamente
	name := fs.String("name", "", "Nombre")
	unit := fs.String("unit", "m", "Unidad") //definiendo la unidad como mega predeterminadamente
	type_ := fs.String("type", "p", "Tipo") //definiendo el tipo de partición como primaria predeterminadamente
	fit := fs.String("fit", "f", "Ajuste") //definiendo el ajuste como first predeterminadamente
	delete := fs.String("fit", "", "Ajuste") //definiendo el valor de delete como vacío predeterminadamente
	add := fs.String("fit", 0, "Ajuste") //definiedo el valor de add como 0 predeterminadamente
	
	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(command, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "driveletter", "name", "type", "delete", "add":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}
	if *driveletter == "" || *name == "" { //Valida que los valores de driveletter y name no estén vacíos ya que son requeridos obligatoriamente
		fmt.Println("Error: Driveletter and name cannot be empty")
		return
	}
	if *unit != "b" && *unit != "k" && *unit != "m" { //Valida que el atributo unit sea igual a bytes, kilo o mega
		fmt.Println("Error: Unit must be b, k or m")
		return
	}
	if *unit == "k" { // Set the size in bytes
		*size = *size * 1024
	} else if *unit == "m" {
		*size = *size * 1024 * 1024
	}
	if *type_ != "p" && *type_ != "e" && *type_ != "l" { //Valida que el atributo type sea igual a p-rimaria, e-xtendida o l-ogica
		fmt.Println("Error: Type must be p, e or l")
		return
	}
	if *fit != "b" && *fit != "w" && *fit != "f" { //Valida que el atributo fit sea igual a best, worst o first
		fmt.Println("Error: Fit must be b, w or f")
		return
	}
	/*if *size <= 0 {//Valida que el atributo size sea mayor a 0
		fmt.Println("Error: Size must be greater than 0")
		return
	}*/

	Commands.FDISK(*size, *driveletter, *name, *unit, *type_, *fit, *delete, *add)

func processMkdiskCommand(command string) {
	fmt.Println("La cadena contiene la palabra 'mkdisk'.")
	size, fit, unit := defineFlags()
	processInput(command, size, fit, unit)
	validateFlags(size, fit, unit)
	if *unit == "k" {
		*size = *size * 1024
	} else {
		*size = *size * 1024 * 1024
	}

	Commands.MKDISK(Files.ObtenerNuevoNombreArchivo(), size, fit, unit)

}

func defineFlags() (*int, *string, *string) {
	size := flag.Int("size", 0, "Tamaño")
	fit := flag.String("fit", "f", "Ajuste")
	unit := flag.String("unit", "m", "Unidad")
	flag.Parse()
	return size, fit, unit
}

func validateFlags(size *int, fit *string, unit *string) {
	if *fit != "b" && *fit != "w" && *fit != "f" {
		fmt.Println("Error: Fit must be b, w or f")
		return
	}
	if *size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}
	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit must be k or m")
		return
	}
	fmt.Println("Size:", *size)
	fmt.Println("Fit:", *fit)
	fmt.Println("Unit:", *unit)
}

func processExecuteCommand(command string) {
	fmt.Println("La cadena contiene la palabra 'EXECUTE'.")
	input, err := subCadena(command, "=", 1) // Verificar si el carácter "-" está presente
	if err != nil {
		return
	}
	fmt.Println("El valor de input es: ", input)
	var lineas []string
	Commands.Execute(input, &lineas)
	for _, linea := range lineas {
		fmt.Println(linea)
		AnalyzeType(linea)
	}
}

func processInput(input string, size *int, fit *string, unit *string) {
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size":
			sizeValue := 0
			fmt.Sscanf(flagValue, "%d", &sizeValue)
			*size = sizeValue
		case "fit":
			flagValue = flagValue[:1]
			flagValue = strings.ToLower(flagValue)
			*fit = flagValue
		case "unit":
			flagValue = strings.ToLower(flagValue)
			*unit = flagValue
		default:
			fmt.Println("Error: Flag not found")
		}
	}
}
func subCadena(command string, value string, sumIndice int) (string, error) {
	indice := strings.Index(command, value)
	var input string
	//input = command[indice:]
	// Verificar si el carácter "-" está presente
	if indice != -1 {
		// Obtener el substring desde la posición hasta el final de la cadena
		input = command[indice+sumIndice:]
	} else {
		fmt.Println("No se ha encontrado el carácter ", value)
	}
	return input, nil
}
