package ui

import (
	"bytes"
	"embed"
	"html/template"
	"io/fs"
)

//go:embed templates
var templateFS embed.FS

func ExpandTemplate(path string, data any) (string, error) {

	content, readErr := fs.ReadFile(templateFS, path)
	if readErr != nil {
		return "", readErr
	}
	t := template.New(path) //.Funcs(funcMap)
	theTemplate, parseErr := t.Parse(string(content))
	if parseErr != nil {
		return "", parseErr
	}
	var buf bytes.Buffer
	execErr := theTemplate.Execute(&buf, data)
	if execErr != nil {
		return "", execErr
	}
	return buf.String(), nil
}
