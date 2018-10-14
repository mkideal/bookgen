package main

import (
	"github.com/mkideal/cli"
)

type checkT struct {
}

var cmdCheck = &cli.Command{
	Name: "check",
	Desc: "check if there is an error in the book",
	Argv: func() interface{} { return new(checkT) },
	Fn: func(ctx *cli.Context) error {
		return book.Run()
	},
}
