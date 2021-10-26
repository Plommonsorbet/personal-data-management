package pdm

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PDM struct {
	Items    []Item
	BasePath string
}

func (p *PDM) Suggest(args []string) {
	for _, item := range p.Items {
		remainder, match := item.PartialMatch(args)
		if match {
			if len(remainder) > 0 {
				fmt.Println(remainder[0])
			}
		}
	}
}

func (p *PDM) Get(path []string) (*Item, error) {
	for _, item := range p.Items {
		if item.IsMatch(path) {
			return &item, nil
		}

	}
	return nil, errors.New("No match for path in PDM.")
}

func (p *PDM) ReadItem(path []string) (string, error) {
	if item, err := p.Get(path); err != nil {
		return "", err
	} else {
		if data, err := item.Read(); err != nil {
			return "", err
		} else {
			return data, nil
		}

	}

}

func LoadPDM(basePath string) (PDM, error) {

	all := []Item{}
	err := filepath.Walk(basePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				rel, err := filepath.Rel(basePath, filepath.Dir(path))
				if err != nil {
					return err
				}
				item := Item{
					Dirs: strings.Split(rel, "/"),
					Name: filepath.Base(path),
					Path: path,
				}
				all = append(all, item)

			}
			return nil
		})

	if err != nil {
		return PDM{}, err
	}
	return PDM{Items: all, BasePath: basePath}, nil
}
