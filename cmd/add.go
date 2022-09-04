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
	fmt.Println("IV is ", iv)
	fmt.Println("CT is ", ct)
	// Generate program tmp/tmp.go
	/*
		1. decrypt bytecode
		2. write bytecode to a tmp file
		3. mark tmp as executable and run
	*/
	prgm := []byte(fmt.Sprintf(`package main

	import (
		"fmt"
		"strconv"
		"strings"
	)
	
	func recvByteSlice(bs string) []byte {
		// Read bytes
		var bb []byte
		for _, ps := range strings.Split(strings.Trim(bs, "[]"), " ") {
			pi, _ := strconv.Atoi(ps)
			bb = append(bb, byte(pi))
		}
		return bb
	}
	
	func main() {
		// s := "WAT"
	
		// Read in iv
		ivStr := fmt.Sprintf("%v")
		iv := recvByteSlice(ivStr)
	
		// Read in iv
		ctStr := fmt.Sprintf("%v")
		ct := recvByteSlice(ctStr)
	
		// Print result
		fmt.Println("IV:", iv)
		fmt.Println("CT:", ct)
	
	}`, iv, ct))

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
	fmt.Println(cmd)
	check(err)

	err = os.Rename("tmp", "../packed")
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
