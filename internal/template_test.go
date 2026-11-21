package internal

import (
	"strings"
	"testing"
)

func renderTemplate(t *testing.T, name string) string {
	t.Helper()

	tmpl, err := loadTemplate(name)
	if err != nil {
		t.Fatalf("loading %s template: %v", name, err)
	}

	var out strings.Builder
	if err := tmpl.Execute(&out, templateData{year: 2024, day: 1}); err != nil {
		t.Fatalf("executing %s template: %v", name, err)
	}

	return out.String()
}

func TestRubyTemplateAssertsInputAnswers(t *testing.T) {
	out := renderTemplate(t, "ruby")

	for _, want := range []string{
		"require 'minitest/autorun'",
		"def read_file(filename)",
		"def test_solve_example",
		"def test_solve_input",
		"def test_solve2_example",
		"def test_solve2_input",
		"solve(read_file('input.txt'))",
		"solve2(read_file('input.txt'))",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("ruby template missing %q", want)
		}
	}
}

func TestRubyTemplateIsSelfContainedAndQuiet(t *testing.T) {
	out := renderTemplate(t, "ruby")

	for _, unwanted := range []string{"require_relative", "puts "} {
		if strings.Contains(out, unwanted) {
			t.Errorf("ruby template should not contain %q", unwanted)
		}
	}
}

func TestPythonTemplateRenders(t *testing.T) {
	if out := renderTemplate(t, "python"); out == "" {
		t.Error("python template rendered empty")
	}
}
