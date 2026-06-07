package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/notomo/gh-create-iteration/createiteration"

	"github.com/urfave/cli/v2"
)

const (
	paramProjectUrl     = "project-url"
	paramIterationField = "field"
	paramCount          = "count"
	paramDuration       = "duration"
	paramStartDate      = "start-date"
	paramTitlePrefix    = "title-prefix"
	paramDryRun         = "dry-run"
	paramLog            = "log"
)

func main() {
	app := &cli.App{
		Name: "gh-create-iteration",
		Action: func(c *cli.Context) error {
			opts := api.ClientOptions{}
			logFilePath := c.String(paramLog)
			if logFilePath != "" {
				f, err := os.Create(logFilePath)
				if err != nil {
					return fmt.Errorf("create log file: %w", err)
				}
				defer f.Close()
				opts.Log = f
				opts.LogVerboseHTTP = true
			}
			gql, err := api.NewGraphQLClient(opts)
			if err != nil {
				return fmt.Errorf("create gql client: %w", err)
			}
			return createiteration.Run(
				gql,
				c.String(paramProjectUrl),
				c.String(paramIterationField),
				c.Int(paramCount),
				c.Int(paramDuration),
				c.String(paramStartDate),
				c.String(paramTitlePrefix),
				c.Bool(paramDryRun),
				os.Stdout,
			)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     paramProjectUrl,
				Value:    "",
				Required: true,
				Usage:    "project url",
			},
			&cli.StringFlag{
				Name:     paramIterationField,
				Value:    "",
				Required: true,
				Usage:    "iteration field name",
			},
			&cli.IntFlag{
				Name:  paramCount,
				Value: 1,
				Usage: "number of iterations to create",
			},
			&cli.IntFlag{
				Name:  paramDuration,
				Value: 0,
				Usage: "duration days per new iteration (0 means inherit the field configuration)",
			},
			&cli.StringFlag{
				Name:  paramStartDate,
				Value: "",
				Usage: "start date (yyyy-mm-dd) of the first new iteration (default: the day after the last existing iteration ends)",
			},
			&cli.StringFlag{
				Name:  paramTitlePrefix,
				Value: "Iteration ",
				Usage: "title prefix for new iterations",
			},
			&cli.BoolFlag{
				Name:  paramDryRun,
				Value: false,
				Usage: "nothing is updated",
			},
			&cli.StringFlag{
				Name:  paramLog,
				Value: "",
				Usage: "log file path",
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
