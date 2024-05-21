package templates

import (
	"html/template"
	"io"
	"strings"
)

type LayoutTemplateProcessor struct{}

func (proc *LayoutTemplateProcessor) ExecTemplate(writer io.Writer, name string, data interface{}) (err error) {
	var sb strings.Builder
	layoutName := ""
	localTemplates := getTemplates()
	localTemplates.Funcs(map[string]interface{}{
		"body":   insertBodyWrapper(&sb),
		"layout": setLayoutWrapper(&layoutName),
	})
	err = localTemplates.ExecuteTemplate(&sb, name, data)
	if err != nil {
		return err
	}

	if layoutName != "" {
		err = localTemplates.ExecuteTemplate(writer, layoutName, data)
	} else {
		_, err = io.WriteString(writer, sb.String())
	}
	return err
}

var getTemplates func() *template.Template

func insertBodyWrapper(body *strings.Builder) func() template.HTML {
	return func() template.HTML {
		return template.HTML(body.String())
	}
}

func setLayoutWrapper(val *string) func(string) string {
	return func(layout string) string {
		*val = layout
		return ""
	}
}
