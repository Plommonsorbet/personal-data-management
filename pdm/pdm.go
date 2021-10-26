package pdm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	if len(i.PathSections()) != len(path) {
		log.Fatal("Path length is not the same")
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
		log.Fatal("Path length is longer than the current path")
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

type PDM struct {
	Items []Item
	BasePath string
}

func (p *PDM) Suggest(args []string) {
	 for _, item := range p.Items {
	 	 remainder, match := item.PartialMatch(args)
	 	 if  match {
		 	 if len(remainder) > 0 {
			 	fmt.Println(remainder[0])
			 }
		 }
	 }
}

// TODO: REMOVE / REPURPOSE - NOT USED
//func (p *PDM) List() [][]string{
//	item_paths := [][]string{}
//	for _, item := range p.Items {
//		item_paths = append(item_paths, item.PathSections())
//	}
//	return item_paths
//}

func (p *PDM) Get(path []string) (*Item, error) {
	for _, item := range p.Items {
		if item.IsMatch(path) {
			return &item, nil
		}

	}
	return nil, errors.New("No match for path in PDM.")
}

func (p *PDM) Read(path []string) (string, error) {
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
		log.Println(err)
		return PDM{}, err
	}
	return PDM{Items: all, BasePath: basePath}, nil
}
