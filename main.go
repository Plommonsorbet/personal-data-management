package main

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/urfave/cli/v2"
	"os"
	"pdm/pdm"
)
//go:embed completions/bash_autocomplete
var bash_complete string
//go:embed completions/zsh_autocomplete
var zsh_complete string

func main() {
	supported_shells := []string{"zsh", "bash"}
	var dataDir string

	app := &cli.App{
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "data",
				Aliases:     []string{"d"},
				Usage:       "Data dir",
				EnvVars:     []string{"PDM_DATA"},
				Value:	     os.Getenv("HOME") + "/.config/pdm",
				Destination: &dataDir,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "completion",
				Usage: "show completion",
				BashComplete: func(c *cli.Context) {
					for _, shell := range supported_shells {
						fmt.Println(shell)
					}
				},
				Action: func(c *cli.Context) error {
					shell := c.Args().First()
					switch shell {
					case "zsh":
						fmt.Println(zsh_complete)
					case "bash":
						fmt.Println(bash_complete)
					default:
						return errors.New(fmt.Sprintf("%s is not one of the supported shells: %s", shell, supported_shells))
					}
					return nil

				},
			},
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
					fmt.Println(dataDir)
					if p, err := pdm.LoadPDM(dataDir); err != nil {
						return err
					} else {
						if data, err := p.ReadItem(c.Args().Slice()); err != nil {
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
						if data, err := p.ReadItem(c.Args().Slice()); err != nil {
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
