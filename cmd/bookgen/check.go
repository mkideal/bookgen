package main

import (
	"github.com/mkideal/cli"
)

type checkT struct {
}

var cmdCheck = &cli.Command{
	Name: "check",
	Desc: "check book",
	Argv: func() interface{} { return new(checkT) },
	Fn: func(ctx *cli.Context) error {
		return book.Run()
	},
}
