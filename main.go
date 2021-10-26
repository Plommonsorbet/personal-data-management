package main

import (
	"pdm/pdm"
	//"fmt"
	"github.com/urfave/cli/v2"
	"os"
	//"log"
	"fmt"
	"github.com/atotto/clipboard"
)

func get_item(dataDir string, args []string) (string, error) {

	if p, err := pdm.LoadPDM(dataDir); err != nil {
		return "", err

	} else {
		if item, err := p.Get(args); err != nil {
			return "", err
		} else {
			if data, err := item.Read(); err != nil {
				return "", err
			} else {
				return data, nil
			}
		}
	}
	return "", nil
}

func main() {
	var dataDir string

	app := &cli.App{
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "data",
				Aliases:     []string{"d"},
				Usage:       "Data dir",
				EnvVars:     []string{"PDM_DATA"},
				Destination: &dataDir,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all items",
				Action: func(c *cli.Context) error {

					p, err := pdm.LoadPDM(dataDir)
					if err != nil {
						return err
					}
					for _, item := range p.Items {
						fmt.Println(item)
					}
					return nil
				},
			},
			{
				Name:    "clip",
				Aliases: []string{"c"},
				Usage:   "clip item",
				BashComplete: func(c *cli.Context) {
					p, _ := pdm.LoadPDM(dataDir)
					p.Suggest(c.Args().Slice())

			},
				Action: func(c *cli.Context) error {
					if p, err := pdm.LoadPDM(dataDir); err != nil {
						return err
					} else {
						if data, err := p.Read(c.Args().Slice()); err != nil {
							return err
						} else {
							clipboard.WriteAll(data)
							fmt.Println("Copy to clipboard")
							return nil
						}
					}

				},
			},
			{
				Name:    "show",
				Aliases: []string{"s"},
				Usage:   "show item",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "clip", Aliases: []string{"c"}},
				},
				BashComplete: func(c *cli.Context) {
					p, _ := pdm.LoadPDM(dataDir)
					p.Suggest(c.Args().Slice())

			},
				Action: func(c *cli.Context) error {
					if p, err := pdm.LoadPDM(dataDir); err != nil {
						return err
					} else {
						if data, err := p.Read(c.Args().Slice()); err != nil {
							return err
						} else {
							fmt.Println(data)
							return nil
						}
					}
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}

}
