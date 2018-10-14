package bookgen

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mkideal/log"
)

var (
	ErrUnexpectedResult = errors.New("unexpected result")
)

type RunResult struct {
	Filename       string
	StartLine      int
	EndLine        int
	CompileError   error
	RunError       error
	Output         string
	ExpectedOutput string
}

func (rr RunResult) Err() error {
	if rr.CompileError != nil {
		return rr.makeError(rr.CompileError, "")
	}
	if rr.RunError != nil {
		return rr.makeResultError(rr.RunError, rr.ExpectedOutput, rr.Output)
	}
	return nil
}

func (rr RunResult) makeError(err error, extra string) error {
	return fmt.Errorf("%s:%d~%d: %v%s", rr.Filename, rr.StartLine, rr.EndLine, err, extra)
}

func (rr RunResult) formatCode(w *bytes.Buffer, code string) {
	lines := strings.Split(code, "\n")
	for j, line := range lines {
		if j != 0 {
			w.WriteByte('\n')
		}
		w.WriteString("\t> ")
		w.WriteString(line)
	}
}

func (rr RunResult) makeResultError(err error, want, got string) error {
	w := new(bytes.Buffer)
	w.WriteString("\n  expected output:\n")
	rr.formatCode(w, want)
	w.WriteString("\n  but got output:\n")
	rr.formatCode(w, got)
	return rr.makeError(err, w.String())
}

type Code struct {
	Filename  string
	StartLine int
	EndLine   int
	Content   []byte `json:"-"`
	Runnable  bool
	Output    string
}

func (code Code) Run() (output string, err error) {
	if !code.Runnable {
		return
	}
	result := RunResult{
		Filename:       code.Filename,
		StartLine:      code.StartLine,
		EndLine:        code.EndLine,
		ExpectedOutput: code.Output,
	}

	// create tmp Go source file
	tmpFilename := filepath.Join(os.TempDir(), RandString(10)+".go")
	var tmpFile *os.File
	tmpFile, err = os.Create(tmpFilename)
	if err != nil {
		return
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFilename)
	if _, err = tmpFile.Write(code.Content); err != nil {
		return
	}

	// run the Go source file
	cmd := exec.Command("go", "run", tmpFilename)
	errWriter := new(bytes.Buffer)
	outWriter := new(bytes.Buffer)
	cmd.Stderr = errWriter
	cmd.Stdout = outWriter
	if err := cmd.Run(); err != nil {
		result.CompileError = errors.New(errWriter.String())
	} else {
		result.Output = outWriter.String()
		if result.Output != code.Output {
			result.RunError = ErrUnexpectedResult
		}
	}
	log.Debug("run %s:%d~%d: \n%s", code.Filename, code.StartLine, code.EndLine, outWriter.String())
	return outWriter.String(), result.Err()
}
