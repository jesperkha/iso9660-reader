package reader

import "strings"

// A File is a struct that describes the data found at the extent location
// of a directory record. No files are read by reading a directory and
// fs.ReadFile() needs to be used instead.
type File struct {
	Name     string
	Size     int
	Position int
	Bytes    []byte
}

// Returns the string version of a file. Only useful for plain text files.
func (file *File) String() string {
	return string(file.Bytes)
}

// ReadFile reads the bytes located at the location described by the parent
// directory record in the given path. Returns a path error if the file is
// not found. Todo: support cross-sector files
func (fs *FileSystem) ReadFile(path string) (file *File, err error) {
	// Get the records for the parent directory
	pathSplit := strings.Split(path, "/")
	parentDir := strings.Join(pathSplit[:len(pathSplit)-1], "/")
	records, err := fs.ReadDirectory(parentDir)
	if err != nil {
		return file, err
	}

	// Find file record
	var record *DirectoryRecord
	for _, r := range records {
		if r.Name == pathSplit[len(pathSplit)-1] {
			record = r
			break
		}
	}

	// Read file
	fileSize, filePos := record.ExtentSize, record.ExtentPos
	content := make([]byte, fileSize)
	_, err = fs.file.ReadAt(content, int64(filePos*blockSize))

	file = &File{
		Name:     record.Name,
		Size:     fileSize,
		Position: filePos,
		Bytes:    content,
	}

	return file, err
}
