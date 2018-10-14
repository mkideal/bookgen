package main

import (
	"fmt"
	"os"

	"github.com/mkideal/bookgen"
	"github.com/mkideal/cli"
	"github.com/mkideal/log"
)

func main() {
	defer log.Uninit(log.InitConsole(log.LvWARN))
	log.NoHeader()

	if err := cli.Root(cmdRoot,
		cli.Tree(cmdCheck),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type rootT struct {
	cli.Helper
	Dir     string      `cli:"d,dir" usage:"root directory of book" dft:"."`
	Verbose cli.Counter `cli:"v" usage:"make the operation more talkative"`
}

var cmdRoot = &cli.Command{
	Desc:   "bookgen used to generate book based on markdown",
	Argv:   func() interface{} { return new(rootT) },
	Global: true,
	OnRootBefore: func(ctx *cli.Context) error {
		argv := ctx.RootArgv().(*rootT)

		switch argv.Verbose.Value() {
		case 0:
			log.SetLevel(log.LvWARN)
		case 1:
			log.SetLevel(log.LvINFO)
		case 2:
			log.SetLevel(log.LvDEBUG)
		default:
			log.SetLevel(log.LvTRACE)
		}

		bookgen.SetRootDir(argv.Dir)
		book = bookgen.Parse(bookgen.Chapters()...)
		log.WithJSON(book).Debug("book")
		return nil
	},
	Fn: func(ctx *cli.Context) error {
		return book.Rewrite()
	},
}
