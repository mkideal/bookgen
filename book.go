package bookgen

type Context struct {
	Book    *Book
	Chapter *Chapter
	Section *Section
}

func (context *Context) WithBook(book *Book) *Context {
	context.Book = book
	return context
}
func (context *Context) WithChapter(chapter *Chapter) *Context {
	context.Chapter = chapter
	return context
}
func (context *Context) WithSection(section *Section) *Context {
	context.Section = section
	return context
}

//---------
// Section
//---------

type Section struct {
	Id   int
	File File
}

func (section Section) Run(errList *ErrorList) {
	section.File.Run(errList)
}

func (section *Section) Rewrite(context Context) error {
	return section.File.Rewrite(*context.WithSection(section))
}

//---------
// Chapter
//---------

type Chapter struct {
	Id       int
	File     File
	Sections []Section
}

func (chapter Chapter) Run(errList *ErrorList) {
	chapter.File.Run(errList)
	for _, sec := range chapter.Sections {
		sec.Run(errList)
	}
}

func (chapter *Chapter) Rewrite(context Context) error {
	context.WithChapter(chapter)
	if err := chapter.File.Rewrite(context); err != nil {
		return err
	}
	for i := 0; i < len(chapter.Sections); i++ {
		if err := chapter.Sections[i].Rewrite(context); err != nil {
			return err
		}
	}
	return nil
}

//------
// Book
//------

type Book struct {
	File     File
	Chapters []Chapter
}

func (book Book) Run() error {
	errList := &ErrorList{}
	for _, ch := range book.Chapters {
		ch.Run(errList)
	}
	return errList.Done()
}

func (book *Book) Rewrite() error {
	var context Context
	context.WithBook(book)
	if err := book.File.Rewrite(context); err != nil {
		return err
	}
	for i := 0; i < len(book.Chapters); i++ {
		if err := book.Chapters[i].Rewrite(context); err != nil {
			return err
		}
	}
	return nil
}
