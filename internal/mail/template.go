package mail

import (
	"bytes"
	"text/template"

	"muni/go-mail/internal/model"
)

// ExecEftTemplate apply info on Template to create html
func ExecEftTemplate(eftinfo model.EftInfo) (bytes.Buffer, error) { 
	eftTemplatePath := "templates/DEFAULT.gohtml"
	// 	eftTemplatePath := "templates/SAMPLE_TEMPLATE.txt" //for testing
	return execTempate(eftTemplatePath, eftinfo)
}

func execTempate(path string, data interface{}) (bytes.Buffer, error) {
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
