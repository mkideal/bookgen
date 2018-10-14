package bookgen

import (
	"testing"
)

func initTestEnv(t *testing.T) {
	*flRootdir = "./"
}

func TestRunCode(t *testing.T) {
	initTestEnv(t)

	code := Code{
		Filename:  "testcode.go",
		StartLine: 0,
		EndLine:   10,
		Runnable:  true,
		Output:    "hello, world",
		Content: []byte(`package main

import "fmt"

func main() {
	fmt.Print("hello, world")
	//panic("errororororororo")
}
`),
	}
	if output, err := code.Run(); err != nil {
		t.Errorf("run code error: %v", err)
	} else {
		t.Logf("output: %s\n", output)
	}
}
