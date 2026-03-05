# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

aoc2md is a Go CLI that initializes Advent of Code puzzle solutions. It fetches the problem description (converted to markdown), input data, and optionally generates a solution template file for a given language (ruby, python).

## Commands

- `make build` — Build binary to `bin/aoc2md`
- `make lint` — Run golangci-lint (config in `.golangci.yml`)
- `make lint-fix` — Run golangci-lint with auto-fix
- `make install` — Install Go dependencies and dev tools (golangci-lint, goreleaser)
- `go vet ./...` — Run Go vet checks
- `gofmt -s -w .` — Format code

There are no tests in this project currently.

## Architecture

Single-command CLI built with `urfave/cli/v2`. Entry point is `main.go` which handles flag parsing and validation. Core logic lives in `internal/action.go`:

- Fetches puzzle HTML from adventofcode.com using a session cookie
- Converts HTML to GitHub-flavored markdown via `html-to-markdown` + `goquery`
- Writes `README.md` (problem description) and `input.txt` to `{year}/day-{DD}/`
- Optionally renders a solution template from `internal/templates/*.tmpl` using Go's `text/template`

Templates are embedded at compile time via `//go:embed`.

## Linting

Uses golangci-lint with a curated set of linters (all disabled by default, explicitly enabled). Import ordering enforced by `gci`: standard library, then `github.com/kevinrobayna`, then `github.com/kevinrobayna/aoc2md`, then third-party. Code formatting uses `gofumpt` with extra rules.

## Release

Uses GoReleaser v2 for releases, distributed via Homebrew (`kevinrobayna/homebrew-tap`). Validate release config with `goreleaser check`.
