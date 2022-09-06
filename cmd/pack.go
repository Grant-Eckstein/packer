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
	"github.com/guumaster/logsymbols"
	"github.com/spf13/cobra"
)

func assertFileExists(filename string) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		printlnFailure("Could not see input file")
		log.Fatal(getError(FileFolderDoesNotExistError, err))
	}
}

func pack(filename string) {
	assertFileExists(filename)

	// Read in file
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		printlnFailure("Failed to read input file")
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
		printlnFailure("Failed to get current working directory")
		log.Fatal(getError(GetCWDError, err))
	}

	// TODO - add tmp dir and have it built there
	tempDir, err := os.MkdirTemp("", "*-temp")
	if err != nil {
		printlnFailure("Failed to create temporary directory")
		log.Fatal(getError(CreateTemporaryFileError, err))
	}
	defer os.RemoveAll(tempDir)
	fmt.Printf("%vCreated %v\n", logsymbols.Success, tempDir)

	// Create temporary file to be built
	tmpFileName := "*-packed.go"
	tmpFile, err := ioutil.TempFile(tempDir, tmpFileName)
	if err != nil {
		printlnFailure("Failed to create temporary file")
		log.Fatal(getError(CreateTemporaryFileError, err))
	}
	defer os.Remove(tmpFile.Name())
	fmt.Printf("%vCreated %v\n", logsymbols.Success, tmpFile.Name())

	// Write template with new values to temporary file
	err = os.WriteFile(tmpFile.Name(), prgm, 0666)
	if err != nil {
		printlnFailure("Failed to write to temporary file")
		log.Fatal(getError(WriteTemporaryFileError, err))
	}
	fmt.Printf("%vWrote to %v\n", logsymbols.Success, tmpFile.Name())

	// Assert that go is installed
	_, err = exec.LookPath("go")
	if err != nil {
		printlnFailure("Could not verify that go is installed")
		log.Fatal(getError(GoNotInstalledError, err))
	}
	fmt.Printf("%vVerified that go is installed\n", logsymbols.Success)

	// Change directories into tempDir
	err = os.Chdir(tempDir)
	if err != nil {
		printlnFailure("Failed Changing Directories into temporary directory")
		log.Fatal(getError(ChangeDirError, err))
	}

	// Initialize module
	cmd := exec.Command("go", "mod", "init", "tmp")
	err = cmd.Run()
	if err != nil {
		printlnFailure("Failed initializing temporary module")
		log.Fatal(getError(GoGetError, err))
	}
	fmt.Printf("%vInitilized temp module\n", logsymbols.Success)

	// Get reqs for temporary file
	cmd = exec.Command("go", "get", "...")
	err = cmd.Run()
	if err != nil {
		printlnFailure("Failed getting new module dependencies")
		log.Fatal(getError(GoGetError, err))
	}
	fmt.Printf("%vGot dependencies for new module\n", logsymbols.Success)

	// Build temporary file
	// TODO - add requirement to cmd to specify target architecture
	cmd = exec.Command("go", "build", tmpFile.Name())
	err = cmd.Run()
	if err != nil {
		printlnFailure("Failed building new module")
		log.Fatal(getError(BuildFailedError, err))
	}
	fmt.Printf("%vBuilt new module\n", logsymbols.Success)

	// Move newly built binary back to cwd
	inFileName := strings.TrimSuffix(tmpFile.Name(), ".go")
	outFileName := path.Join(cwd, "packed")
	err = os.Rename(inFileName, outFileName)
	if err != nil {
		printlnFailure("Failed moving new executable")
		log.Fatal(getError(CopyTempFileToOutput, err))
	}
	printlnSuccess("Moved new executable")
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
