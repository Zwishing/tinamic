package query

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"sync"
	"text/template"
)

// cache SQL/HTML templates so repeated filesystem reads are not required
var globalTemplates map[string](*template.Template) = make(map[string](*template.Template))
var globalTemplatesMutex = &sync.Mutex{}

func GetSQLTemplate(name string, tmpl string) *template.Template {
	tp, ok := globalTemplates[name]
	if ok {
		return tp
	}
	t := template.New(name)
	tp, err := t.Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}
	globalTemplatesMutex.Lock()
	globalTemplates[name] = tp
	globalTemplatesMutex.Unlock()
	return tp
}

func RenderSQLTemplate(name string, tmpl string, data interface{}) (string, error) {
	var buf bytes.Buffer
	t := GetSQLTemplate(name, tmpl)
	err := t.Execute(&buf, data)
	if err != nil {
		return string(buf.Bytes()), err
	}
	sql := string(buf.Bytes())
	log.Debug(sql)
	return sql, nil
}
