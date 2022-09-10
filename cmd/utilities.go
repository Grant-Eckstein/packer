package cmd

import (
	"bytes"
	"compress/flate"
	"fmt"
	"github.com/guumaster/logsymbols"
)

func printlnSuccess(s string) {
	fmt.Printf("%v%v\n", logsymbols.Success, s)
}

func printlnFailure(s string) {
	fmt.Printf("%v%v\n", logsymbols.Error, s)
}

func compressBytes(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := flate.NewWriter(&b, 9)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
