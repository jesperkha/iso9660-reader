package reader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

var (
	ErrFileTooSmall = errors.New("file size is too small to be a valid ISO image")
	ErrNotValidFieldType = errors.New("field type is invalid")
	ErrNotValidISOFile = errors.New("a standard identifier was not found, file is assumed invalid")
)

const (
	kilobyte   = 1024
	blockSize  = 2048

	// Size of reserved boot information before primary descriptor
	bootSize = 32 * kilobyte
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

// Field in a sector table; a byte array interval. Most fields are either
// a/d strings for names and identifiers, or unsigned integers. Fields start
// at an offset of bytes to not look through empty slots.
type TableField struct {
	Name 	 string
	Offset   int
	Length   int
	DataType int
}

// Todo: find a better way of storing the values in a struct or smth
var PrimaryTable = []TableField{
	// {"System identifier", 8, 8, STR_A},
	// {"Volume identifier", 40, 32, STR_D},
	// {"Volume space size", 80, 8, INT32_LSB_MSB},
	// {"Number of disks", 120, 4, INT16_LSB_MSB},
	// {"Volume sequence number", 124, 4, INT16_LSB_MSB},
	// {"Logical block size", 128, 4, INT16_LSB_MSB},
	// {"Path table size", 132, 8, INT32_LSB_MSB},
	// {"Path table location", 140, 4, INT32},
	// {"Optional table location", 144, 4, INT32},
	// {"Volume set identifier", 190, 128, STR_D},
	// {"Publisher identifier", 318, 128, STR_A},
	// {"Data preparer identifier", 446, 128, STR_A},
	// {"Application identifier", 574, 128, STR_A},
	// {"Copyright file", 702, 38, STR_D},
	// {"Abstract file", 740, 36, STR_D},
	// {"Bibliografic file", 776, 37, STR_D},
	{"Root address", 156, 34, DIRECTORY_RECORD},
}

// Todo: make convert to descriptor struct
// Reads the primary descriptor of the image disk. Does not modify the file.
// If the length is less than (32KB + block_size) an error is returned. If
// the sector identifier is not found, an error is returned.
func ReadPrimaryDescriptor(f *os.File) (desc map[string]interface{}, err error) {
	// Check if file size is big enough to avoid checking for EOF
	if s, _ := f.Stat(); s.Size() < bootSize + blockSize {
		return desc, ErrFileTooSmall
	}
	
	// Get sector interval, size of 2048 bytes
	sector := make([]byte, blockSize)
	if _, err = f.ReadAt(sector, 32*kilobyte); err != nil {
		return desc, err
	}

	// Check for ISO standard identifier found in the first 5 bytes after the
	// sector type. Is always 'CD001', else file is assumed invalid.
	stdIdentifer := string(sector[1:6])
	if stdIdentifer != "CD001" {
		return desc, ErrNotValidISOFile
	}

	data := map[string]interface{}{}
	for _, field := range PrimaryTable {
		// Get byte interval for field
		start := field.Offset
		end := field.Offset + field.Length
		// LSB_MSB has duplicate values for both little- and big endian.
		// Here we are only read the little endian (first) values.
		if field.DataType == INT32_LSB_MSB || field.DataType == INT16_LSB_MSB {
			end -= field.Length / 2
		}

		interval := sector[start:end]
		var value interface{}

		switch field.DataType {
		case STR_A, STR_D:
			value = string(interval)
		case INT32_LSB_MSB, INT32:
			value = binary.LittleEndian.Uint32(interval)
		case INT16_LSB_MSB, INT16:
			value = binary.LittleEndian.Uint16(interval)
		case INT8:
			value = uint8(interval[0])
		case DEC_TIME:
			fmt.Println("time not implemented")
		case DIRECTORY_RECORD:
			fmt.Println("dir not implemented")
		default:
			return desc, ErrNotValidFieldType
		}

		data[field.Name] = value
	}

	return data, err
}