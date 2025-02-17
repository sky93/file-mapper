package main

import (
	"log"
	"os"

	"github.com/sky93/file-mapper/internal/listing"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "Project Mapper",
		Usage:    "A simple tool to map your project tree and file contents.",
		Version:  "v0.0.1",
		Commands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Root path to scan",
				Value:   ".", // default
			},
			&cli.StringFlag{
				Name:    "include",
				Aliases: []string{"i"},
				Usage:   "Comma-separated file patterns to include (e.g. '*.go,*.txt')",
			},
			&cli.StringFlag{
				Name:    "exclude",
				Aliases: []string{"e"},
				Usage:   "Comma-separated directories/files to exclude (e.g. '.git,.idea,.env')",
				// By default, we are ignoring hidden dirs. If you want to
				// always exclude e.g. ".git,.idea,.env" add a default Value here.
			},
			&cli.BoolFlag{
				Name:    "git",
				Aliases: []string{"g"},
				Usage:   "Only list Git-tracked files",
			},
			&cli.BoolFlag{
				Name:    "content",
				Aliases: []string{"c"},
				Usage:   "Include file content (for text files)",
			},
			&cli.BoolFlag{
				Name:  "separate-content",
				Usage: "If set, print the tree first, then print all file contents afterward",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file path (if not provided, prints to stdout)",
			},
			&cli.BoolFlag{
				Name:  "flat",
				Usage: "Show files in a flat list instead of the default tree",
			},
			&cli.BoolFlag{
				Name:  "line-numbers",
				Usage: "Show line numbers for file content",
			},
			&cli.BoolFlag{
				Name:  "header-footer",
				Usage: "Show '----- CONTENT START -----' and '----- CONTENT END -----' markers",
				Value: true, // default is to show them
			},
		},
		Action: func(ctx *cli.Context) error {
			cfg := &listing.Config{
				RootPath:        ctx.String("path"),
				Include:         ctx.String("include"),
				Exclude:         ctx.String("exclude"),
				GitTrackedOnly:  ctx.Bool("git"),
				ShowTree:        !ctx.Bool("flat"), // default is tree
				ShowContent:     ctx.Bool("content"),
				SeparateContent: ctx.Bool("separate-content"),
				Output:          ctx.String("output"),

				ShowLineNumbers:   ctx.Bool("line-numbers"),
				ShowHeaderFooters: ctx.Bool("header-footer"),
			}

			result, err := listing.Run(cfg)
			if err != nil {
				return err
			}

			if cfg.Output != "" {
				f, err := os.Create(cfg.Output)
				if err != nil {
					return err
				}
				defer f.Close()
				if _, err = f.WriteString(result); err != nil {
					return err
				}
				log.Printf("Output written to %s\n", cfg.Output)
			} else {
				os.Stdout.WriteString(result)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
