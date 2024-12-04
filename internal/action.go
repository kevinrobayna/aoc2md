package internal

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/PuerkitoBio/goquery"
)

type Session string

type Template string

const (
	None Template = "none"
	Ruby Template = "ruby"
)

type Args struct {
	Session  Session
	Language Template
	Day      int
	Year     int
}

type templateData struct {
	year int
	day  int
}

//go:embed templates/*.tmpl
var templateFS embed.FS

func GenerateTemplate(args Args) error {
	year := args.Year
	day := args.Day
	session := args.Session

	path := generatePath(year, day)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		slog.Error("Error creating folder", "err", err)
		return nil
	}

	slog.Info("Fetching Problem Description")
	description, err := fetchDescription(year, day, session)
	if err != nil {
		slog.Error("Whopsy, the elf was not able to fetch the problem for you", "err", err)
		return nil
	}
	err = writeToFile(filepath.Join(path, "README.md"), description)
	if err != nil {
		slog.Error("Unable to write description into README.md", "err", err)
		return nil
	}

	if _, err := os.Stat(filepath.Join(path, "input.txt")); errors.Is(err, os.ErrNotExist) {
		slog.Info("Fetching Problem Input")
		input, err := fetchInput(year, day, session)
		if err != nil {
			slog.Error("Whopsy, the elf was not able to fetch the problem for you", "err", err)
			return nil
		}
		err = writeToFile(filepath.Join(path, "input.txt"), input)
		if err != nil {
			slog.Error("Unable to write input into input.txt", "err", err)
			return nil
		}
	} else {
		slog.Warn("Skipping downloading input, file already exist")
	}

	if args.Language == None {
		slog.Warn("Skipping code template generation")
	} else {
		tmpl, err := loadTemplate(string(args.Language))
		if err != nil {
			slog.Error("Error loading template", "err", err)
			return err
		}

		if _, err := os.Stat(filepath.Join(path, generateSolutionName(args.Language))); errors.Is(err, os.ErrNotExist) {
			slog.Info("Creating template for problem")

			outputFile, err := os.Create(filepath.Join(path, generateSolutionName(args.Language)))
			if err != nil {
				slog.Error("Error creating file", "err", err)
				return err
			}
			defer outputFile.Close()

			data := templateData{
				year: args.Year,
				day:  args.Day,
			}
			err = tmpl.Execute(outputFile, data)
			if err != nil {
				slog.Error("Error executing template", "err", err)
				return err
			}
		} else {
			slog.Warn("Skipping code template generation, file already exist")
		}
	}

	return nil
}

func generatePath(year, day int) string {
	return filepath.Join(fmt.Sprint(year), fmt.Sprintf("day-%02d", day))
}

func generateSolutionName(lang Template) string {
	if lang == Ruby {
		return "main.rb"
	}
	slog.Error("unknown language", "lang", lang)
	return "main.txt"
}

func prepareRequest(url string, session Session) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("Unable to prepare request", "err", err)
		return nil, err
	}

	cookie := &http.Cookie{
		Name:  "session",
		Value: string(session),
	}
	req.AddCookie(cookie)

	// This is for us to tell AoC maitainer where the requests are coming from.
	// We aren't required to do this but we want to be good citizens
	req.Header.Add("User-Agent", "github.com/kevinrobayna/aoc2md")

	return req, nil
}

func fetchInput(year, day int, session Session) (string, error) {
	client := &http.Client{}
	req, err := prepareRequest(fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day), session)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Unable to make request", "err", err)
		return "", err
	}
	defer resp.Body.Close()

	input, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Unable to read input from response", "err", err)
		return "", err
	}

	return string(input), err
}

func fetchDescription(year, day int, session Session) (string, error) {
	client := &http.Client{}
	req, err := prepareRequest(fmt.Sprintf("https://adventofcode.com/%d/day/%d", year, day), session)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Unable to make request", "err", err)
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		slog.Error("Unable to parse HTML response from AOC", "err", err)
		return "", err
	}
	// Find the content of the <article> element with class "day-desc"
	converter := md.NewConverter(md.DomainFromURL("https://adventofcode.com"), true, nil)
	// Use the `GitHubFlavored` plugin from the `plugin` package.
	converter.Use(plugin.GitHubFlavored())
	markdown := converter.Convert(doc.Find(".day-desc"))
	return markdown, nil
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

func loadTemplate(name string) (*template.Template, error) {
	templateContent, err := templateFS.ReadFile("templates/" + name + ".tmpl")
	if err != nil {
		return nil, err
	}

	// Parse the template from the content
	return template.New(name).Parse(string(templateContent))
}
