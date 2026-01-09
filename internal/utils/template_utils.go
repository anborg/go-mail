package utils

import (
	"bytes"
	"text/template"
)

// RenderTemplate parses a template file and executes it with the provided data.
func RenderTemplate(path string, data interface{}) (bytes.Buffer, error) {
	var bytesHtml bytes.Buffer
	tpl, err := template.ParseFiles(path)
	if err != nil {
		return bytesHtml, err
	}

	if err := tpl.Execute(&bytesHtml, data); err != nil {
		return bytesHtml, err
	}
	return bytesHtml, nil
}
