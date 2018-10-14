package main

import (
	"fmt"
	"os"

	"github.com/mkideal/bookgen"
	"github.com/mkideal/cli"
	"github.com/mkideal/log"
	"github.com/mkideal/log/logger"
)

func main() {
	if err := cli.Root(cmdRoot,
		cli.Tree(cmdCheck),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type rootT struct {
	cli.Helper
	Dir string       `cli:"d,dir" usage:"root directory" dft:"."`
	Log logger.Level `cli:"log" usage:"log level: trace/debug/info/warn/error/fatal"`
}

var cmdRoot = &cli.Command{
	Desc:   "bookgen used to generate book based on markdown",
	Argv:   func() interface{} { return new(rootT) },
	Global: true,
	OnRootBefore: func(ctx *cli.Context) error {
		argv := new(rootT)
		if err := ctx.GetArgvAt(argv, 0); err != nil {
			return err
		}
		log.SetLevel(argv.Log)
		bookgen.SetRootDir(argv.Dir)
		book = bookgen.Parse(bookgen.Chapters()...)
		return nil
	},
	Fn: func(ctx *cli.Context) error {
		return book.Rewrite()
	},
}
