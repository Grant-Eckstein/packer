package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/Grant-Eckstein/everglade"
	"github.com/spf13/cobra"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func assertFileExists(filename string) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		errorText := fmt.Sprintf("File '%s' does not exist", filename)
		log.Fatal(errorText)
	}
}

func pack(filename string) {
	assertFileExists(filename)

	// Read in file
	fileBytes, err := os.ReadFile(filename)
	check(err)

	// Generate keys, salt, and IV
	eg := everglade.New()

	// Encrypt bytes
	iv, ct := eg.EncryptCBC(fileBytes)
	exp := eg.Export()
	// Generate program tmp/tmp.go
	/*
		1. decrypt bytecode [DONE]
		2. write bytecode to a tmp file
		3. mark tmp as executable and run
	*/
	// TODO - replace this with exporting eg.Export() JSON
	prgm := []byte(fmt.Sprintf(`package main

	import (
		"fmt"
		"io/ioutil"
		"log"
		"os"
		"os/exec"
		"strconv"
		"strings"
	
		"github.com/Grant-Eckstein/everglade"
	)
	
	func recvByteSlice(bs string) []byte {
		var bb []byte
		for _, ps := range strings.Split(strings.Trim(bs, "[]"), " ") {
			pi, _ := strconv.Atoi(ps)
			bb = append(bb, byte(pi))
		}
		return bb
	}
	
	func main() {
	
		// Read in iv
		ivStr := fmt.Sprintf("%v")
		iv := recvByteSlice(ivStr)
	
		// Read in iv
		ctStr := fmt.Sprintf("%v")
		ct := recvByteSlice(ctStr)
	
		// Read in everglade export
		exStr := fmt.Sprintf("%v")
		exp := recvByteSlice(exStr)
	
		obj := everglade.Import(exp)
	
		data := obj.DecryptCBC(iv, ct)
	
		// Create temp file	
		file, err := ioutil.TempFile(".", ".*")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name())
	
		// Write to file
		err = os.WriteFile(file.Name(), data, 0777)
		if err != nil {
			log.Fatal(err)
		}

		err = os.Chmod(file.Name(), 0777)
		if err != nil {
			log.Fatal(err)
		}
	
		err = exec.Command(file.Name()).Run()
		
		if err != nil {
			log.Fatal(err)
		}
	}
	`, iv, ct, exp))

	// Write new golang program tmp/tmp.go
	_, err = os.Stat("tmp")
	if os.IsExist(err) {
		log.Fatal("tmp directory already exists, please remove this and rerun.")
	}
	err = os.Mkdir("tmp", 0777)
	check(err)

	filePath := path.Join("tmp", "packed.go")
	err = os.WriteFile(filePath, prgm, 0666)
	check(err)

	// run go Build tmp.go
	err = os.Chdir("tmp")
	check(err)

	_, err = exec.LookPath("go")
	check(err)

	cmd := exec.Command("go", "build")
	err = cmd.Run()
	check(err)

	newFolder := path.Join("..", "packed")
	err = os.Rename("tmp", newFolder)
	check(err)

	// run go Build tmp.go
	err = os.Chdir("..")
	check(err)

	err = os.RemoveAll("tmp")
	check(err)

}

// packCmd represents the add command
var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "pack an executables bytecode",
	Long:  `pack an executables bytecode`,
	Run: func(cmd *cobra.Command, files []string) {
		// For each specified file, pack
		for _, file := range files {
			// Assert that file exists
			pack(file)

			// Pack file
		}
	},
}

func init() {
	rootCmd.AddCommand(packCmd)
}
