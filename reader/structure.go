package reader

// Todo: finish building structure with dates and stuff
type DescriptorHeader struct {
	Type       int
	Identifier string
	Version    int
}

type DirectoryRecord struct {
	Length    int

	FileFlags int
	Filename  string
	// Date 
	
}

// The primary volume descriptor a collection of information on the disk file.
// Mainly, the root directory record is used to traverse the file tree.
type PrimaryDescriptor struct {
	Header        *DescriptorHeader
	RootDirectory *DirectoryRecord

	// Size of file disk volume and block
	// Block size should be 2048 bytes according to the ISO 9660 standard
	VolumeSize int
	BlockSize  int

	PathTableSize     int
	PathTableLocation int

	// Additional system information
	Publisher        string
	DataPreparer     string
	VolumeIdentifier string
}
