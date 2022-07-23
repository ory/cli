package comments

import "path/filepath"

// a file format that we know about, represented as its file extension
type FileType string

// indicates whether the given list of FileTypes contains the given FileType
func ContainsFileType(fileTypes []FileType, fileType FileType) bool {
	for _, ft := range fileTypes {
		if ft == fileType {
			return true
		}
	}
	return false
}

// provides the extension of the given filename
func GetFileType(filename string) FileType {
	ext := filepath.Ext(filename)
	if len(ext) == 0 {
		return ""
	}
	return FileType(ext[1:])
}

// indicates whether it is possible to add comments to the file with the given name
func Supports(filename string) bool {
	filetype := GetFileType(filename)
	_, ok := commentFormats[filetype]
	return ok
}
