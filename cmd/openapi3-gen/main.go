package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/skygeario/openapi3-gen/internal"
)

var baseDir string
var outputFile string

func init() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&baseDir, "dir", workDir, "project base directory")
	flag.StringVar(&outputFile, "output", "", "output OpenAPI specification file")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] <patterns...>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	patterns := flag.Args()

	if len(patterns) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	err := internal.Run(baseDir, patterns, outputFile)
	if err != nil {
		panic(err)
	}
}
