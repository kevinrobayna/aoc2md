package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"slices"
	"time"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/kevinrobayna/aoc2md/internal"
)

var (
	appName = "aoc2md"
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	cli.VersionPrinter = func(_ *cli.Context) {
		_, _ = fmt.Printf(
			"Application: %s\nVersion: %s\nCommit: %s\nGo Version: %v\nGo OS/Arch: %v/%v\nBuilt at: %v\n",
			appName, version, commit, runtime.Version(), runtime.GOOS, runtime.GOARCH, date,
		)
	}
	app := &cli.App{
		Name:        appName,
		Description: "CLI to help initialise the solution for Advent of Code 🎄. Solving it though it on you 😉",
		Version:     version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "session",
				Aliases:  []string{"s"},
				Usage:    "Session granted by adventofcode.com when you log in",
				Required: false,
				EnvVars:  []string{"AOC_SESSION"},
			},
			&cli.IntFlag{
				Name:        "day",
				Aliases:     []string{"d"},
				Usage:       "Day you're trying to initialize",
				Required:    false,
				Value:       time.Now().Day(),
				DefaultText: fmt.Sprint(time.Now().Day()),
			},
			&cli.IntFlag{
				Name:        "year",
				Aliases:     []string{"y"},
				Usage:       "Year you're trying to initialize",
				Required:    false,
				Value:       time.Now().Year(),
				DefaultText: fmt.Sprint(time.Now().Year()),
			},
			&cli.StringFlag{
				Name:     "lang",
				Aliases:  []string{"l"},
				Usage:    "Programming language to generate the solution for",
				Required: false,
				Value:    "none",
				Action: func(ctx *cli.Context, v string) error {
					if v == "" {
						return nil
					}
					allowedValues := []string{string(internal.Ruby)}
					if slices.Contains(allowedValues, v) {
						return nil
					}
					return fmt.Errorf("Invalid language provided '%v', we currently support: %v", v, allowedValues)
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			handler := log.NewWithOptions(os.Stderr, log.Options{
				ReportCaller:    false,
				ReportTimestamp: false,
				Prefix:          "🎄aoc2md🎄",
			})
			session := internal.Session(ctx.String("session"))
			day := ctx.Int("day")
			year := ctx.Int("year")
			lang := ctx.String("lang")
			slog.SetDefault(slog.New(handler).With("day", day, "year", year, "lang", lang))
			if session == "" {
				slog.Error("The session is required so that you can download your input and part2. see --help")
				return nil
			}
			if year < 2015 {
				slog.Error("Advent Of Code started on 2015, there are no problems before then")
				return nil
			}
			if year > time.Now().Year() {
				slog.Error("Sorry but I can't see in the future")
				return nil
			}
			if year == time.Now().Year() && day > time.Now().Day() {
				slog.Error("Be patient! Don't try to get the problem from a future day")
				return nil
			}
			if day < 1 || day > 25 {
				slog.Error("Advent Of Code problems are only available from the 1st of December to the 25th")
				return nil
			}
			args := internal.Args{
				Session:  session,
				Day:      day,
				Year:     year,
				Language: internal.Template(lang),
			}
			return internal.GenerateTemplate(args)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
