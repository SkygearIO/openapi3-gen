package internal

import (
	"fmt"
	"strings"

	"github.com/skygeario/openapi3-gen/pkg/processor"
	"github.com/skygeario/openapi3-gen/pkg/scanner"
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

	spec, errs := psr.End()
	if len(errs) > 0 {
		return Error{errs}
	}

	// TODO: output spec in yaml
	fmt.Println(spec)

	return nil
}
