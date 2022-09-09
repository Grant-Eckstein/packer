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

	/*** Generate new build ***/
	// Load template
	prgmFile, err := os.ReadFile("cmd/template")
	if err != nil {
		printlnFailure("Failed to read in template")
		log.Fatal(getError(ReadTemplateFileError, err))
	}
	prgm := string(prgmFile)

	// Insert IV into template
	ivStr := fmt.Sprintf("%v", iv)
	prgm = strings.Replace(prgm, "INSERT_IV", ivStr, 1)

	// Insert CT into template
	ctStr := fmt.Sprintf("%v", ct)
	prgm = strings.Replace(prgm, "INSERT_CT", ctStr, 1)

	// Insert everglade export into template
	expStr := fmt.Sprintf("%v", exp)
	prgm = strings.Replace(prgm, "INSERT_EXP", expStr, 1)

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
	err = os.WriteFile(tmpFile.Name(), []byte(prgm), 0666)
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
	// env GOOS=target-OS GOARCH=target-architecture
	goosBuildString := fmt.Sprintf("GOOS=%v", Goos)
	goarchBuildString := fmt.Sprintf("GOARCH=%v", Goarch)
	cmd = exec.Command("env", goosBuildString, goarchBuildString, "go", "build", tmpFile.Name())
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

var Goos string
var Goarch string

func init() {
	rootCmd.AddCommand(packCmd)
	// define required local flag
	packCmd.Flags().StringVarP(&Goos, "goos", "o", "", "Set build os")
	packCmd.MarkFlagRequired("goos")

	packCmd.Flags().StringVarP(&Goarch, "goarch", "a", "", "Set build architecture")
	packCmd.MarkFlagRequired("goarch")
}
