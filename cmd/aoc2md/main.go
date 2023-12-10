package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

var (
	appName = "aoc2md"
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

type Session string

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
				Usage:    "Session granted by adventofcode.com when you log in",
				Required: false,
				EnvVars:  []string{"AOC_SESSION_ID"},
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
		},
		Action: func(ctx *cli.Context) error {
			handler := log.NewWithOptions(os.Stderr, log.Options{
				ReportCaller:    false,
				ReportTimestamp: false,
				Prefix:          "🎄aoc2md🎄",
			})
			session := Session(ctx.String("session_id"))
			if session == "" {
				slog.Error("The session is required so that you can download your input and part2. see --help")
				return nil
			}
			day := ctx.Int("day")
			year := ctx.Int("year")
			slog.SetDefault(slog.New(handler).With("day", day, "year", year))
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

			path := filepath.Join(fmt.Sprint(year), fmt.Sprintf("day-%01d", day))
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				slog.Error("Error creating folder", "err", err)
				return nil
			}

			slog.Info("Fetching Problem Description")
			description, err := FetchDescription(year, day, session)
			if err != nil {
				slog.Error("Whopsy, the elf was not able to fetch the problem for you", "err", err)
			}
			err = writeToFile(filepath.Join(path, "README.md"), description)
			if err != nil {
				slog.Error("Unable to write description into README.md", "err", err)
			}

			if _, err := os.Stat(filepath.Join(path, "input.txt")); errors.Is(err, os.ErrNotExist) {
				slog.Info("Fetching Problem Input")
				input, err := FetchInput(year, day, session)
				if err != nil {
					slog.Error("Whopsy, the elf was not able to fetch the problem for you", "err", err)
				}
				err = writeToFile(filepath.Join(path, "input.txt"), input)
				if err != nil {
					slog.Error("Unable to write input into input.txt", "err", err)
				}
			} else {
				slog.Warn("Skipping downloading input")
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func writeToFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
