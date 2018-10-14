package bookgen

import (
	"io"
)

type BlockType int

const (
	TextBlock BlockType = iota
	CodeBlock
	ReplBlock
)

type Block struct {
	Type    BlockType
	Tag     string `json:".omitempty"`
	Prefix  string `json:"-"`
	Suffix  string `json:"-"`
	Content string `json:"-"`
}

func (block *Block) Rewrite(context Context, w io.Writer) error {
	if block.Type == ReplBlock {
		switch block.Tag {
		case "@toc_of_book":
			// TODO
		case "@toc_of_chapter":
			// TODO
		}
	}

	return nil
}
