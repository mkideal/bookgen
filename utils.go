package bookgen

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	charset = []byte("ABCDEFGHIJKLMNOPQRSATUVWXYZabcdefghijklmnopqrsatuvwxyz")

	ErrChapterOutOfRange = errors.New("chapter out of range: must be in interval [1,99]")
	ErrSectionOutOfRange = errors.New("section out of range: must be in interval [1,99]")

	README = "README.md"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var rootdir string

func SetRootDir(dir string) {
	rootdir = dir
}

func RootDir() string {
	return rootdir
}

func RandString(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		c := charset[rand.Intn(len(charset))]
		s += string(c)
	}
	return s
}

func CheckChapterNumber(chapter int) error {
	if chapter <= 0 || chapter >= 100 {
		return ErrChapterOutOfRange
	}
	return nil
}

func CheckSectionNumber(section int) error {
	if section <= 0 || section >= 100 {
		return ErrSectionOutOfRange
	}
	return nil
}

func ChapterName(chapter int) string {
	return fmt.Sprintf("ch%02d", chapter)
}

func ChapterDir(chapter int) string {
	return filepath.Join(RootDir(), ChapterName(chapter))
}

func ChapterFilepath(chapter int) string {
	chapterName := ChapterName(chapter)
	return filepath.Join(RootDir(), chapterName, README)
}

func SectionFilenamePrefix(chapter, section int) string {
	return fmt.Sprintf("%2d-", section)
}

func SectionFilepath(chapter, section int, sectionName string) string {
	chapterName := ChapterName(chapter)
	return filepath.Join(RootDir(), chapterName, fmt.Sprintf("%2d-%s.md", section, sectionName))
}

func ChapterExist(chapter int) bool {
	_, err := os.Stat(ChapterDir(chapter))
	return err == nil
}

func SectionExist(chapter, section int) bool {
	chapterName := ChapterName(chapter)
	pattner := filepath.Join(RootDir(), chapterName, fmt.Sprintf("-%2d-*.md", section))
	matchs, err := filepath.Glob(pattner)
	return err == nil && len(matchs) > 0
}

func CreateChapter(chapter int) (exist bool, err error) {
	if err := CheckChapterNumber(chapter); err != nil {
		return false, err
	}
	if ChapterExist(chapter) {
		return true, nil
	}
	return false, os.MkdirAll(ChapterDir(chapter), 0755)
}

type ErrorList struct {
	errs []error
}

func (e *ErrorList) Done() error {
	if len(e.errs) == 0 {
		return nil
	}
	return e
}

func (e *ErrorList) Push(err error) {
	if err != nil {
		e.errs = append(e.errs, err)
	}
}

func (e *ErrorList) Error() string {
	b := new(bytes.Buffer)
	for i, err := range e.errs {
		if i != 0 {
			b.WriteByte('\n')
		}
		b.WriteString(err.Error())
	}
	return b.String()
}

func parseChapterFromDirname(dirname string) (chapter int, ok bool) {
	if len(dirname) == 4 && dirname[:2] == "ch" {
		v, err := strconv.ParseInt(dirname[2:], 10, 64)
		if err == nil && v > 0 && v < 100 {
			return int(v), true
		}
	}
	return 0, false
}

var sectionFilenameRegexp = regexp.MustCompile("^([0-9]{2})-.*\\.md$")

func parseSectionFromFilename(filename string) (section int, ok bool) {
	if filename == README {
		return 0, true
	}
	ret := sectionFilenameRegexp.FindStringSubmatch(filename)
	if len(ret) == 2 {
		sectionId, err := strconv.Atoi(ret[1])
		if err == nil && sectionId > 0 && sectionId < 100 {
			return int(sectionId), true
		}
	}
	return 0, false
}

func Chapters() []int {
	chapters := []int{}
	filepath.Walk(RootDir(), func(path string, info os.FileInfo, err error) error {
		if info == nil || err != nil {
			return filepath.SkipDir
		}
		if info.IsDir() {
			chapter, ok := parseChapterFromDirname(info.Name())
			if ok {
				chapters = append(chapters, chapter)
			}
			if RootDir() != path {
				return filepath.SkipDir
			}
		}
		return nil
	})
	return chapters
}

func Sections(chapter int) (sections []Section) {
	chapterDir := ChapterDir(chapter)
	filepath.Walk(chapterDir, func(path string, info os.FileInfo, err error) error {
		if info == nil || err != nil {
			return filepath.SkipDir
		}
		if info.IsDir() && path != chapterDir {
			return filepath.SkipDir
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		id, ok := parseSectionFromFilename(info.Name())
		if ok {
			sections = append(sections, Section{
				Id: id,
				File: File{
					Filename: filepath.Join(RootDir(), info.Name()),
				},
			})
		}
		return nil
	})
	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Id < sections[j].Id
	})
	return
}
