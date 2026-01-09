package eftnotify

import (
	"bytes"
	"muni/go-mail/internal/utils"
)

// GenerateEmailBody apply info on Template to create html
func GenerateEmailBody(eftinfo EftInfo) (bytes.Buffer, error) { 
	eftTemplatePath := "templates/DEFAULT.gohtml"
	// 	eftTemplatePath := "templates/SAMPLE_TEMPLATE.txt" //for testing
	return utils.RenderTemplate(eftTemplatePath, eftinfo)
}

