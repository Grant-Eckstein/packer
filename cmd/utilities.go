package cmd

import (
	"fmt"
	"github.com/guumaster/logsymbols"
)

func printlnSuccess(s string) {
	fmt.Printf("%v%v\n", logsymbols.Success, s)
}

func printlnFailure(s string) {
	fmt.Printf("%v%v\n", logsymbols.Error, s)
}
