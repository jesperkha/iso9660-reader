package reader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotADirectory = errors.New("sector does not contain any records")
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

// FindDirectory traverses forward through the file system to find the specified
// record. If no record with the given name is found an error is returned. See
// DirectoryRecord.ReadFile() for getting the file contents.
// 
// Path must be formatted as "path/to/file.ext" or "path/to/dir". Relative or absolute
// paths are not accepted and the target will not be found.
func (fs *FileSystem) FindDirectory(path string) (dirs []DirectoryRecord, err error) {
	// Get position of root directory. First two entries are the default ./ and ../
	// and both have a length of 34 bytes.
	rootPos := fs.Descriptor.RootDirLocation
	pathSplit := strings.Split(path, "/")
	location := rootPos // keep track of current sector

	for _, dirname := range pathSplit {
		// Check records to match dirname and continue to subfolder
		dirs, err = fs.ReadDirectory(location)
		if err != nil {
			return dirs, err
		}

		var target *DirectoryRecord
		for _, r := range dirs {
			if r.Name == dirname {
				target = &r
				break
			}
		}

		if target == nil {
			return dirs, fmt.Errorf(ErrDirectoryNotFound.Error(), dirname)
		}

		location = target.ExtentPos
	}

	// Return final read of the target directory
	return fs.ReadDirectory(location)
}

// ReadDirectory reads a list of directory records from the given sector number (from 0).
// Records spanning across multiple sectors are allowed. Returns a list of records including
// the ./ and ../ entries. Todo: allow multi-sector records
func (fs *FileSystem) ReadDirectory(location int) (dirs []DirectoryRecord, err error) {
	sector := make([]byte, blockSize)
	_, err = fs.File.ReadAt(sector, int64(location * blockSize))
	if err != nil {
		return dirs, err
	}

	// Verify there are records in this sector. Not foolproof but good enough.
	if int(sector[0]) != 34 || int(sector[34]) != 34 {
		return dirs, ErrNotADirectory
	}

	dirs = []DirectoryRecord{}
	index := 0
	for index < blockSize {
		length := int(sector[index])
		if length == 0 { // end of records
			break
		}
		
		// Get record interval and increment index
		interval := sector[index:index+length]
		index += length

		// Read needed field values
		record := DirectoryRecord{
			Size: length,
			ExtentPos: int(binary.LittleEndian.Uint32(interval[2:6])),
			ExtentSize: int(binary.LittleEndian.Uint32(interval[10:14])),
			Flag: int(interval[25]),
		}

		// Set names of first two entries to . and .. for convenience
		name := string(interval[33:33+int(interval[32])])
		switch interval[33] {
		case 0:
			name = "."
		case 1:
			name = ".."
		}

		record.Name = name
		dirs = append(dirs, record)
	}

	return dirs, err
}

// ReadFile writes the bytes of the file described by the given record. Returns an error
// if the record does not describe a file. Todo: support cross sector files
func (fs *FileSystem) ReadFile(dest []byte, record DirectoryRecord) (err error) {

	return err
}