package Mbr

type MBR struct {
	Tamano        int32
	FechaCreacion int64
	DiskSignature int32
	// EspacioReservado [fileSize - 16]byte // Espacio restante en la estructura
}
