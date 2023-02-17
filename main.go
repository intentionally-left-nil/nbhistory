package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "nbhistory",
		Commands: []*cli.Command{
			{
				Name:      "save",
				Usage:     "Add the current state of a notebook to git history",
				ArgsUsage: "Filenames to save",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "message", Aliases: []string{"m"}, Usage: "Save description", Required: true},
				},
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() == 0 {
						return errors.New("missing filename")
					}
					for _, filename := range ctx.Args().Slice() {
						err := save(filename, ctx.String("message"))
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			{
				Name:      "load",
				Usage:     "Restore to the last saved version of the notebook",
				ArgsUsage: "Filenames to load",
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() == 0 {
						return errors.New("missing filename")
					}
					for _, filename := range ctx.Args().Slice() {
						err := load(filename)
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
			{
				Name:      "revert",
				Usage:     "Forget the most recently saved notebook, and go back to the previous one",
				ArgsUsage: "Filenames to revert",
				Action: func(ctx *cli.Context) error {
					if ctx.Args().Len() == 0 {
						return errors.New("missing filename")
					}
					for _, filename := range ctx.Args().Slice() {
						err := revert(filename)
						if err != nil {
							return err
						}
						err = load(filename)
						if err != nil {
							return err
						}
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
