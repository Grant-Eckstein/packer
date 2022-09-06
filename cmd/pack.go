package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/Grant-Eckstein/everglade"
	"github.com/spf13/cobra"
)

func assertFileExists(filename string) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		log.Fatal(getError(FileFolderDoesNotExistError, err))
	}
}

func pack(filename string) {
	assertFileExists(filename)

	// Read in file
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(getError(ReadFileError, err))
	}

	// Generate keys, salt, and IV
	eg := everglade.New()

	// Encrypt bytes
	iv, ct := eg.EncryptCBC(fileBytes)
	exp := eg.Export()

	// Generate new build
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

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(getError(GetCWDError, err))
	}

	// TODO - add tmp dir and have it built there
	tempDir, err := os.MkdirTemp("", "*-temp")
	if err != nil {
		log.Fatal(getError(CreateTemporaryFileError, err))
	}
	defer os.RemoveAll(tempDir)
	fmt.Printf("Created temporary dir - %v\n", tempDir)

	// Create temporary file to be built
	tmpFileName := "*-packed.go"
	tmpFile, err := ioutil.TempFile(tempDir, tmpFileName)
	if err != nil {
		fmt.Println(err)
		log.Fatal(getError(CreateTemporaryFileError, err))
	}
	defer os.Remove(tmpFile.Name())
	fmt.Printf("Created temporary file - %v\n", tmpFile.Name())

	// Write template with new values to temporary file
	err = os.WriteFile(tmpFile.Name(), prgm, 0666)
	if err != nil {
		fmt.Println(err)
		log.Fatal(getError(WriteTemporaryFileError, err))
	}
	fmt.Printf("Wrote to temporary file - %v\n", tmpFile.Name())

	// Assert that go is installed
	_, err = exec.LookPath("go")
	if err != nil {
		fmt.Println(err)
		log.Fatal(getError(GoNotInstalledError, err))
	}
	fmt.Println("Verified that go is installed")

	// Change directories into tempDir
	err = os.Chdir(tempDir)
	if err != nil {
		log.Fatal(getError(ChangeDirError, err))
	}

	// Initialize module
	cmd := exec.Command("go", "mod", "init", "tmp")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Output is", string(out))
		log.Fatal(getError(GoGetError, err))
	}
	fmt.Println("Initilized temp module")

	// Get reqs for temporary file
	cmd = exec.Command("go", "get", "...")
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Output is", string(out))
		log.Fatal(getError(GoGetError, err))
	}
	fmt.Println("Got dependencies for new module")

	// Build temporary file
	cmd = exec.Command("go", "build", tmpFile.Name())
	err = cmd.Run()
	if err != nil {
		log.Fatal(getError(BuildFailedError, err))
	}
	fmt.Println("Built new module")

	// Move newly built binary back to cwd
	inFileName := strings.TrimSuffix(tmpFile.Name(), ".go")
	outFileName := path.Join(cwd, "packed")
	err = os.Rename(inFileName, outFileName)
	if err != nil {
		log.Fatal(getError(CopyTempFileToOutput, err))
	}
	fmt.Println("Moved new executable")

	// // Change directories back to cwd
	// err = os.Chdir(cwd)
	// if err != nil {
	// 	log.Fatal(getError(ChangeDirError, err))
	// }
	// fmt.Println("Changed back to cwd")

	// // Copy tmp file to regular file with target name
	// outFileName := fmt.Sprintf("packed-%v", filename)
	// err = copyFile(tmpFileName, outFileName)
	// if err != nil {
	// 	log.Fatal(getError(CopyTempFileToOutput, err))
	// }

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
		}
	},
}

func init() {
	rootCmd.AddCommand(packCmd)
}
