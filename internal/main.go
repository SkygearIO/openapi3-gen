package internal

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/skygeario/openapi3-gen/pkg/processor"
	"github.com/skygeario/openapi3-gen/pkg/scanner"
	"gopkg.in/yaml.v2"
)

type Error struct {
	inner []error
}

func (err Error) Error() string {
	lines := []string{"failed to process: "}
	for _, err := range err.inner {
		lines = append(lines, err.Error())
	}
	return strings.Join(lines, "\n")
}

func Run(baseDir string, patterns []string, outputFile string) error {
	psr := processor.New()
	scn := scanner.New(psr.Process)

	err := scn.Scan(baseDir, patterns)
	if err != nil {
		return err
	}

	oapi, errs := psr.End()
	if len(errs) > 0 {
		return Error{errs}
	}

	oapiData, err := yaml.Marshal(oapi)
	if err != nil {
		return err
	}

	if outputFile != "" {
		err = ioutil.WriteFile(outputFile, oapiData, 0644)
	} else {
		_, err = fmt.Print(string(oapiData))
	}

	return err
}
