package cmd

import (
	"errors"
	"fmt"
)

const (
	EmptyError                  = "Error reported with empty message"
	GetCWDError                 = "Could not get CWD"
	FileFolderDoesNotExistError = "File or folder does not exist"
	ReadFileError               = "Could not read file"
	CreateTemporaryFileError    = "Could not create temporary file"
	WriteTemporaryFileError     = "Could not write temporary file"
	GoNotInstalledError         = "go executable not installed on path"
	GoGetError                  = "Could not run go get"
	BuildFailedError            = "Could not build go file"
	ChangeDirError              = "Unable to change directories"
	WriteFileError              = "Could not write file"
	CopyTempFileToOutput        = "Could not copy temporary built executable"
)

func getError(s string, e error) error {
	// When empty, use manual label
	if s != "" {
		return errors.New(s)
	}
	// Include friendly message and origional error
	msg := fmt.Sprintf("%v: %v", s, e)
	return errors.New(msg)
}
