package Models

import "fmt"

type MBR struct {
	Tamano            int32
	FechaCreacion     int64
	DiskSignature     int32
	Fit               [1]byte
	Partitions        [4]Partition
	PartitionExtended [1]PartitionExtended
	// EspacioReservado [fileSize - 16]byte // Espacio restante en la estructura
}

func PrintMBR(data MBR) {
	fmt.Println(fmt.Sprintf("CreationDate: %s, fit: %s, size: %d", string(data.FechaCreacion), string(data.Fit[:]), data.Tamano))
	for i := 0; i < 4; i++ {
		PrintPartition(data.Partitions[i])
	}
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

type PartitionExtended struct {
	Status          [1]byte
	Type            [1]byte
	Fit             [1]byte
	Start           int32
	Size            int32
	Name            [16]byte
	Correlative     int32
	Id              [4]byte
	PartitionLogica [20]PartitionLogica
}

func PrintExtendedPartition(data PartitionExtended) {
	fmt.Println(fmt.Sprintf("Extended Partition - Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s", string(data.Name[:]), string(data.Type[:]), data.Start, data.Size, string(data.Status[:]), string(data.Id[:])))
	for i := 0; i < 20; i++ {
		PrintPartitionLogica(data.PartitionLogica[i])
	}
}

func PrintPartitionLogica(data PartitionLogica) {
	fmt.Println(fmt.Sprintf("Logical Partition - Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s", string(data.Name[:]), string(data.Type[:]), data.Start, data.Size, string(data.Status[:]), string(data.Id[:])))
}

func PrintMBRWithExtended(data MBR) {
	PrintMBR(data)
	PrintExtendedPartition(data.PartitionExtended[0])
}

type PartitionLogica struct {
	Status      [1]byte
	Type        [1]byte
	Fit         [1]byte
	Start       int32
	Size        int32
	Name        [16]byte
	Correlative int32
	Id          [4]byte
}

func PrintPartition(data Partition) {
	fmt.Println(fmt.Sprintf("Name: %s, type: %s, start: %d, size: %d, status: %s, id: %s", string(data.Name[:]), string(data.Type[:]), data.Start, data.Size, string(data.Status[:]), string(data.Id[:])))
}

// Extended Boot Record (EBR)
type EBR struct {
	part_mount [1]byte  //Indica si la partición está montada o no
	part_fit   [1]byte  // Tipo de ajuste de la partición. Tendrá los valores B (Best), F (First) o W (worst)
	part_start int32    // Indicae en qué byte del disco inicia la partición
	part_s     int32    //Contiene el tamaño total de la partición en bytes.
	part_next  int32    // Byte  en el que está el próximo EBR. -1 si no hay siguiente
	part_name  [16]byte //[16] Nombre de la partición
}

//  =============================================================

type Superblock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [17]byte
	S_umtime            [17]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_fist_ino          int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

//  =============================================================

type Inode struct {
	I_uid   int32
	I_gid   int32
	I_size  int32
	I_atime [17]byte
	I_ctime [17]byte
	I_mtime [17]byte
	I_block [15]int32
	I_type  [1]byte
	I_perm  [3]byte
}

//  =============================================================

type Fileblock struct {
	B_content [64]byte
}

//  =============================================================

type Content struct {
	B_name  [12]byte
	B_inodo int32
}

type Folderblock struct {
	B_content [4]Content
}

//  =============================================================

type Pointerblock struct {
	B_pointers [16]int32
}

//  =============================================================

type Content_J struct {
	Operation [10]byte
	Path      [100]byte
	Content   [100]byte
	Date      [17]byte
}

type Journaling struct {
	Size      int32
	Ultimo    int32
	Contenido [50]Content_J
}
