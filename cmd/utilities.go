package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/guumaster/logsymbols"
)

func printlnSuccess(s string) {
	fmt.Printf("%v%v\n", logsymbols.Success, s)
}

func printlnFailure(s string) {
	fmt.Printf("%v%v\n", logsymbols.Error, s)
}

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
