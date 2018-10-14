package bookgen

import (
	"bytes"
	"io/ioutil"
	"os"
)

//------
// File
//------

type File struct {
	Filename string
	Mode     os.FileMode
	Title    string
	CodeList []Code
	Blocks   []Block
}

func (file File) Run(errList *ErrorList) {
	for _, code := range file.CodeList {
		if _, err := code.Run(); err != nil {
			errList.Push(err)
		}
	}
}

func (file *File) Rewrite(context Context) error {
	var buf bytes.Buffer
	for i := range file.Blocks {
		if err := file.Blocks[i].Rewrite(context, &buf); err != nil {
			return err
		}
	}
	mode := file.Mode
	if mode == 0 {
		mode = 0666
	}
	return ioutil.WriteFile(file.Filename, buf.Bytes(), mode)
}
