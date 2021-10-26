package pdm

import (
	"fmt"
	"os"
	"strings"
)

type Item struct {
	Dirs []string
	Name string
	Path string
}

func (i *Item) PathSections() []string {
	return append(i.Dirs, i.Name)
}

func (i *Item) Read() (string, error) {
	dat, err := os.ReadFile(i.Path)
	if err != nil {
		return "", err
	}
	text := string(dat)
	return text, nil

}

func (i Item) String() string {
	return fmt.Sprintf(strings.Join(i.PathSections(), " "))
}

func (i *Item) IsMatch(path []string) bool {
	if len(i.PathSections()) < len(path) {
		return false
	}

	for i, itemPathSection := range i.PathSections() {
		if itemPathSection != path[i] {
			return false
		}

	}
	return true
}

func (i *Item) PartialMatch(args []string) ([]string, bool) {
	if len(i.PathSections()) < len(args) {
		return []string{}, false
	}

	for j, arg := range args {
		if i.PathSections()[j] != arg {
			return []string{}, false
		}

	}

	remainder := i.PathSections()[len(args):]

	return remainder, true
}
