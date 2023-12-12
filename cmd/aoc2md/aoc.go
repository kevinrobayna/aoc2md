package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
	"github.com/PuerkitoBio/goquery"
)

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

func FetchInput(year, day int, session Session) (string, error) {
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

func FetchDescription(year, day int, session Session) (string, error) {
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
	if err != nil {
		slog.Error("Was that HTML? Uhm something went wrong", "err", err)
		return "", err
	}
	return markdown, nil
}
