package reader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrDirectoryNotFound = errors.New("could not find directory '%s'")
)

// A directory record contains information on the file or folder it describes.
// Records cannot spill over into another sector which means, if the amount
// of records exceeds a sector, padding of 0-bytes is added to reach the
// desired size of 2048 bytes per sector.
//
// All diectories have two constant record enties. One is the "./" path and
// the other is the "../" path, in that order. They have the names "\0" and
// "\1" (bytes 0 and 1).
type DirectoryRecord struct {
	// Size of this record in bytes
	Size int
	Flag int

	// Name of directory or file.
	// Filenames are formatted as FILENAME.EXT;VERSION
	Name string

	// Size of extent. Either directory record or file contents
	ExtentSize int
	ExtentPos  int
}

// ReadDirectory traverses forward through the file system to find the specified
// record. If no record with the given name is found an error is returned. Only
// returns the record. See DirectoryRecord.ReadFile() for getting the file contents.
// 
// Path must be formatted as "path/to/file.ext" or "path/to/dir". Relative or absolute
// paths are not accepted and the target will not be found.
func (fs *FileSystem) ReadDirectory(path string) (dirs []*DirectoryRecord, err error) {
	// Get position of root directory. First two entries are the default ./ and ../
	// and both have a length of 34 bytes.
	rootPos := fs.Descriptor.RootDirLocation
	pathSplit := strings.Split(path, "/")

	for _, dirname := range pathSplit {
		// Load entire sector into memory and
		sector := make([]byte, blockSize)
		_, err := fs.File.ReadAt(sector, int64(rootPos))
		if err != nil {
			return dirs, err
		}

		// Read records in current the directory and locate target directory
		dirs = []*DirectoryRecord{}
		pointer := 0

		var target *DirectoryRecord
		for {
			length := int(sector[pointer])
			if length == 0 { // End of records
				break
			}
			
			// Get record interval and increment pointer
			interval := sector[pointer:pointer+length]
			pointer += length

			// Read needed field values
			record := &DirectoryRecord{
				Size: length,
				ExtentPos: int(binary.LittleEndian.Uint32(interval[2:6])),
				ExtentSize: int(binary.LittleEndian.Uint32(interval[10:14])),
				Flag: int(interval[25]),
			}

			// Check if name matches target
			recordName := string(interval[33:33+int(interval[32])])
			if recordName == dirname {
				target = record
			}
			
			record.Name = recordName
			dirs = append(dirs, record)
		}

		if target == nil {
			return dirs, fmt.Errorf(ErrDirectoryNotFound.Error(), path)
		}

		// Todo: start new search
		// Todo: make a separate function?
	}

	return dirs, err
}

// ReadFile writes the bytes of the file described by the given record. Returns an error
// if the record does not describe a file. Todo: support cross sector files
func (fs *FileSystem) ReadFile(dest []byte, record *DirectoryRecord) (err error) {

	return err
}