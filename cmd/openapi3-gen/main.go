package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/skygeario/openapi3-gen/internal"
)

var outputFile string

func init() {
	flag.StringVar(&outputFile, "output", "", "output OpenAPI specification file")
	flag.StringVar(&outputFile, "o", "", "output OpenAPI specification file (shorthand)")
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

	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = internal.Run(workDir, patterns, outputFile)
	if err != nil {
		panic(err)
	}
}
