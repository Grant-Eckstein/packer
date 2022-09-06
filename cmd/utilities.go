package cmd

import (
	"io/ioutil"
)

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
