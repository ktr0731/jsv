package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

func showError(err error) {
	fmt.Fprintf(os.Stderr, "\x1b[31mjsv: %s\x1b[0m\n", err)
}

func run(args []string) int {
	if len(args) < 2 {
		showError(errors.New("Usage: jsv <input> <schema ...>"))
		return 1
	}

	// TODO: Windows
	aPath, err := filepath.Abs(args[0])
	if err != nil {
		showError(err)
		return 1
	}

	is := gojsonschema.NewReferenceLoader("file://" + aPath)

	for _, uri := range args[1:] {
		vaPath, err := filepath.Abs(uri)
		if err != nil {
			showError(err)
			return 1
		}

		vs := gojsonschema.NewReferenceLoader("file://" + vaPath)
		res, err := gojsonschema.Validate(vs, is)
		if err != nil {
			showError(err)
			return 1
		}

		if !res.Valid() {
			errMsg := fmt.Sprintf("validation failed: schema: %s\n", uri)
			for _, err := range res.Errors() {
				errMsg += err.String() + "\n"
			}
			showError(errors.New(errMsg))
			return 1
		}
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
