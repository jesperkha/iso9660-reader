package reader

import (
	"encoding/binary"
	"errors"
	"os"
	"strings"
)

var (
	ErrNotValidFieldType = errors.New("field type is invalid")
	ErrNotValidISOFile   = errors.New("a standard identifier was not found, file is assumed invalid")
)

const (
	// Data types for sector values. LSB indicates little endian
	// values, MSB is big endian. Types such as INT32_LSB_MSB are
	// used for sector parts that are stored as both little- and
	// big endian, where the first half of the part is the little
	// endian byte array. This application only uses the little
	// endian values for consistancy, but there is no difference.
	INT8 = iota
	INT16
	INT32
	INT16_LSB_MSB
	INT32_LSB_MSB

	// Different types of string groups. a-characters are A through
	// Z (including lower case) in addition to digits and some symbols.
	// d-characters are only upper case A-Z, digits, and underscore.
	// Filenames traditionally use d-characters including one period "."
	// and a semicolon ";" to indicate FILENAME.EXT;
	STR_A
	STR_D

	// Date/time format used in primary descriptor. Represents time
	// and date using ASCII digits. Size of 17 bytes.
	DEC_TIME

	DIRECTORY_RECORD
)

type FileSystem struct {
	Descriptor *PrimaryDescriptor
	file       *os.File

	// Keeps track of visted paths and stores them to avoid loading directories
	// more than once. Map key is the position of the first record.
	cachedDirs map[int][]*DirectoryRecord
}

// The primary volume descriptor a collection of information on the disk file.
// Mainly, the root directory record is used to traverse the file tree.
type PrimaryDescriptor struct {
	Type int

	// Size of entire disk
	VolumeSize int

	PathTableSize     int
	PathTableLocation int
	RootDirLocation   int

	// Additional system information
	Publisher        string
	DataPreparer     string
	VolumeIdentifier string
}

// ReadDisk reads the primary descriptor from the disk. Returns a file system
// handler for reading the rest of the files and directories.
func ReadDisk(name string) (fs *FileSystem, err error) {
	file, err := os.Open(name)
	if err != nil {
		return fs, err
	}

	fs = &FileSystem{file: file, cachedDirs: make(map[int][]*DirectoryRecord)}

	// Check for 'CD0001' identifier that is always found in the first sector
	identifier := fs.seekFieldValue(1, 5, STR_A).(string)
	if identifier != "CD001" {
		return fs, ErrNotValidISOFile
	}

	// Read the primary descriptor.
	fs.Descriptor = &PrimaryDescriptor{
		Type:             fs.seekFieldValue(0, 1, INT8).(int),
		VolumeSize:       fs.seekFieldValue(80, 8, INT32_LSB_MSB).(int),
		VolumeIdentifier: fs.seekFieldValue(40, 32, STR_D).(string),
		Publisher:        fs.seekFieldValue(318, 128, STR_A).(string),
		RootDirLocation:  fs.seekFieldValue(158, 8, INT16_LSB_MSB).(int),
	}

	return fs, err
}

// Gets the field value located at the target offset. Returns value based on given
// data	type. If EOF is found function returns nil.
func (fs *FileSystem) seekFieldValue(offset int64, length int, datatype int) interface{} {
	// LSB_MSB has duplicate values for both little- and big endian.
	// Here we are only read the little endian (first) values.
	if datatype == INT16_LSB_MSB || datatype == INT32_LSB_MSB {
		length /= 2
	}

	interval := make([]byte, length)
	// Get interval at default offset (after the reserved system area)
	n, _ := fs.file.ReadAt(interval, sysAreaOffset+offset)
	if n != length { // EOF
		return nil
	}

	switch datatype {
	case STR_A, STR_D:
		return strings.TrimSpace(string(interval))
	case INT32_LSB_MSB, INT32:
		return int(binary.LittleEndian.Uint32(interval))
	case INT16_LSB_MSB, INT16:
		return int(binary.LittleEndian.Uint16(interval))
	case INT8:
		return int(interval[0])
	}

	return nil
}

// Closes file in use. Cannot perform read operations after close
func (fs *FileSystem) Close() {
	fs.file.Close()
}
