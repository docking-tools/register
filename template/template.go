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
    templates map[string][]*config.ConfigTemplate
    url string // uri path
}

func NewTemplate(config *config.ConfigFile) (api.RegistryAdapter, error) {
    u, err := url.ParseRequestURI(config.RegisterUrl)
    
    if err!=nil {
        return nil, err
    }

	if len(config.Templates)==0 {
	    return nil, errors.New("No template found.")
	}
	parseTemplates(config.Templates)
	return &TemplateRegistry{templates: config.Templates, url: u.String()}, nil
}

func (r *TemplateRegistry) Size() int {
    return len(r.templates)
}

func (r *TemplateRegistry) Ping() error {
    return nil
}

func parseTemplates(confTmpl map[string][]*config.ConfigTemplate) {

	for _,confList := range confTmpl {
	    for _,conf := range confList {
    	    conf.SetTmpl(template.Must(template.New(conf.Name).Parse(conf.Template)))
            log.Println("New template: ", conf)
	    }
	}
}

func (r *TemplateRegistry) RunTemplate(status string, service *api.Service) error {
    tmpls := []*config.ConfigTemplate {}
    tmpls = append(tmpls, r.templates[status]...)
    tmpls = append(tmpls, r.templates["ALL"]...) 

    log.Println("%v register for service %v", status, service)

	for _, tmpl := range tmpls {
	    query,err :=executeTemplates(tmpl, service)
	    if err != nil {
            return err
        }
        err = exectureQuery(r.url, query, tmpl.HttpCmd)  
        if err != nil {
            return err
        }        
	}

	return nil    
}

func executeTemplates(conf *config.ConfigTemplate, service *api.Service) (string,error) {

    bufQuery := &bytes.Buffer {}
    // Execute the template with the service as the data item
    bufQuery.Reset()

    err := conf.Tmpl().Execute(bufQuery, service)
    if err != nil {
        return "", err
    }
      

    return bufQuery.String(), nil
}

func exectureQuery(url string, tmpl string, httpCmd string) error {
        client := &http.Client{}    
        querys := strings.Split(tmpl,"\n")

    for _,query := range querys {        
        queryTab := strings.SplitN(query, " ", 2)
        path := queryTab[0]
        value := ""
        if (len(queryTab)==2) {
            value = queryTab[1]
        }

        if len(path) > 0 {
            request, err := http.NewRequest(httpCmd, url+path, strings.NewReader(value))
            request.ContentLength = int64(len(value))
            log.Println("Query: "+url+path+" "+ value)
            
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
                log.Print("response "+string(contents))
            }
        }
    }
    return nil
}