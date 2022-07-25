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

// provides the extension of the given filepath
func GetFileType(filePath string) FileType {
	ext := filepath.Ext(filePath)
	if len(ext) > 0 {
		ext = ext[1:]
	}
	return FileType(ext)
}

// indicates whether it is possible to add comments to the file with the given name
func SupportsFile(filePath string) bool {
	filetype := GetFileType(filePath)
	_, ok := commentFormats[filetype]
	return ok
}
