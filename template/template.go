package template

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"net/http"
	"strings"
	"text/template"
	"github.com/docking-tools/register/api"
	"github.com/docking-tools/register/config"
)

// Usage:
// @TODO

type TemplateRegistry struct {
	api.RegistryAdapter
	templates  map[string][]*config.ConfigTemplate
	httpHeader map[string]*template.Template
	url        string // uri path
}

var funcs = map[string]interface{}{
	"env": env,
	"convertGraphTopath": convertGraphTopath,
	"listPathfromGraph": listPathfromGraph,
}

func NewTemplate(config *config.ConfigFile) (api.RegistryAdapter, error) {
	u, err := url.ParseRequestURI(config.RegisterUrl)

	if err != nil {
		return nil, err
	}

	if len(config.Templates) == 0 {
		return nil, errors.New("No template found.")
	}

	return &TemplateRegistry{templates: parseTemplates(config.Templates), httpHeader: parseMap(config.HttpHeaders), url: u.String()}, nil
}

func (r *TemplateRegistry) Size() int {
	return len(r.templates)
}

func (r *TemplateRegistry) Ping() error {
	return nil
}

func parseMap(data map[string]string) map[string]*template.Template {
	result := make(map[string]*template.Template)
	for key, value := range data {
		result[key] = template.Must(template.New(key).Funcs(funcs).Parse(value))
	}
	return result
}

func parseTemplates(confTmpl map[string][]*config.ConfigTemplate) map[string][]*config.ConfigTemplate {

	for key, confList := range confTmpl {

		for _, conf := range confList {
			conf.SetTmpl(template.Must(template.New(conf.Name).Funcs(funcs).Parse(conf.Template)))
		}
		if (strings.Contains(key, ",")) {
			keys := strings.Split(key, ",")
			for _, newKey := range keys {
				confTmpl[strings.ToUpper(newKey)] = append(confTmpl[strings.ToUpper(newKey)], confList...)
			}
			delete(confTmpl, key)
		}
	}

	return confTmpl
}

func (r *TemplateRegistry) RunTemplate(status string, object interface{}) error {
	tmpls := []*config.ConfigTemplate{}
	tmpls = append(tmpls, r.templates[strings.ToUpper(status)]...)
	tmpls = append(tmpls, r.templates["ALL"]...)
	if (len(tmpls) < 1) {
		log.Printf("no template found for event %s %v", strings.ToUpper(status), r.templates)
	}

	// calcul httpHeader for all query
	headers, err := executeHttpHeaders(r.httpHeader, object)
	if err != nil {
		return err
	}
	for _, tmpl := range tmpls {
		query, err := executeTemplates(tmpl, object)
		if err != nil {
			return err
		}
		err = exectureQuery(r.url, query, tmpl.HttpCmd, headers)
		if err != nil {
			return err
		}
	}

	return nil
}

func executeTemplates(conf *config.ConfigTemplate, object interface{}) (string, error) {
	bufQuery := &bytes.Buffer{}
	// Execute the template with the object as the data item
	bufQuery.Reset()

	err := conf.Tmpl().Execute(bufQuery, object)
	if err != nil {
		return "", err
	}

	return bufQuery.String(), nil
}

func executeHttpHeaders(data map[string]*template.Template, object interface{}) (map[string]string, error) {
	result := make(map[string]string)
	for key, value := range data {
		bufQuery := &bytes.Buffer{}
		err := value.Execute(bufQuery, object)
		if err != nil {
			log.Fatal("Error on execute httpHeader template: %s ", key)
		} else {
			result[key] = bufQuery.String()
		}
	}
	return result, nil
}

func exectureQuery(url string, tmpl string, httpCmd string, httpHeaders map[string]string) error {
	client := &http.Client{}
	querys := strings.Split(tmpl, "\n")

	for _, query := range querys {
		queryTab := strings.SplitN(query, " ", 2)
		path := queryTab[0]
		value := ""
		if (len(queryTab) == 2) {
			value = queryTab[1]
		}

		if len(path) > 0 {
			request, err := http.NewRequest(httpCmd, url + path, strings.NewReader(value))
			request.ContentLength = int64(len(value))
			for key, value := range httpHeaders {
				request.Header.Add(key, value)
			}
			response, err := client.Do(request)
			if (err != nil) {
				log.Fatal(err)
				return err
			}
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
				return err
			} else {
				log.Printf("Query: %s / response %s ", url + path, string(contents))
			}
		}
	}
	return nil
}