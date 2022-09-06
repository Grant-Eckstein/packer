package cmd

import (
	"io/ioutil"
)

// Create temporary directory with specific permission int - Might not need

// Create temporary file and write specified bytes with specific permission int

func copyFile(sourceFile string, destinationFile string) error {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return getError(ReadFileError, err)
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		return getError(WriteFileError, err)
	}
	return nil
}
