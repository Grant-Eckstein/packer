package cmd

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestAssertFileExists(t *testing.T) {
	// Create temp file
	file, err := ioutil.TempFile(".", ".*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	err = assertFileExists(file.Name())
	if err != nil {
		t.Fatal(err)
	}
}
