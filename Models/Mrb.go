package Mbr

type MBR struct {
	Tamano        int32
	FechaCreacion int64
	DiskSignature int32
	Fit           [1]byte
	Partitions    [4]Partition
	// EspacioReservado [fileSize - 16]byte // Espacio restante en la estructura
}
type Partition struct {
	Status      [1]byte
	Type        [1]byte
	Fit         [1]byte
	Start       int32
	Size        int32
	Name        [16]byte
	Correlative int32
	Id          [4]byte
}
