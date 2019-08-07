package main

import (
	"flag"
	"fmt"
	"os"
)

var baseDir string
var outputFile string

func init() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&baseDir, "dir", workDir, "project base directory")
	flag.StringVar(&outputFile, "output", "", "output OpenAPI specification file (stdout if empty)")
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

	err := run(baseDir, patterns, outputFile)
	if err != nil {
		panic(err)
	}
}
