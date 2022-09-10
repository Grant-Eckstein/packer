package cmd

import (
	"errors"
	"fmt"
	"github.com/guumaster/logsymbols"
)

const (
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
	ReadTemplateFileError       = "Could not read template file"
	CompressFileError           = "Could not compress input file"
)

func getError(s string, e error) error {
	if s == "" {
		return errors.New(s)
	}
	msg := fmt.Sprintf("%v%v: %v", logsymbols.Error, s, e)
	return errors.New(msg)
}
