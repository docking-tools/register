package template

import (
    "log"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/docking-tools/register/config"
    "github.com/docking-tools/register/api"
)

func TestNewEmptyClientUrl(t *testing.T) {
    config := config.NewConfigFile("test")
    templ , err := NewTemplate(&config)
    assert.Nil(t, templ)
    assert.NotNil(t, err)
    
}

func TestNewNoTemplate(t *testing.T) {

    config := config.NewConfigFile("test")
    templ, err :=NewTemplate(&config)
    assert.NotNil(t, err)
    assert.Nil(t, templ)
    //assert.Equal(t, &TemplateRegistry{templates: templates, url: "http://test:2375/"}, templ)
}

func TestNewWithGoodtemplate(t *testing.T) {

    // Run init
    conf := config.NewConfigFile("test")    
    conf.RegisterUrl = "http://localhost"
    templates:= conf.Templates
    template := config.ConfigTemplate {
        Name:       "TEST",
        HttpCmd:    "PUT",
        Template:      "/my/query data exemple",
    }
    templates["ALL"] = append(templates["ALL"], &template)
    
    templ, err :=NewTemplate( &conf)
    // Check
    assert.Nil(t, err)
    assert.NotNil(t, templ)
    
}

func TestExecuteTemplate(t *testing.T) {
    conf := config.NewConfigFile("")
    templates :=conf.Templates 
    template := config.ConfigTemplate {
        Name:       "TEST",
        HttpCmd:    "et",
        Template:      "/my/query/{{.ID}} data exemple {{.Name}}",
    }
    templates["ALL"] = append(templates["ALL"], &template)
    parseTemplates(templates)
    
     log.Println("Generate config file %v ", config.ConfigDir(), conf.Filename())

    conf.Save()


    service:= api.Service {
        ID: "idddd",
    	Name: "container A",
    	Port: 8080,
    	IP:    "0.0.0.0",
    	Tags:  make([]string,0),
    	Attrs: make(map[string]string),
    }
    
    assert.NotNil(t, templates["ALL"])
    assert.NotNil(t, templates["ALL"][0])
    
    query,err := executeTemplates(templates["ALL"][0], &service)
    log.Println("Executed template: %v %v", query)
    assert.Nil(t, err)
    assert.NotNil(t, query)
    assert.Equal(t, "/my/query/idddd data exemple container A", query)
}

func TestStructMultiQuery(t *testing.T) {
    conf := config.NewConfigFile("")
    templates :=conf.Templates 
    template := config.ConfigTemplate {
        Name:       "TEST",
        HttpCmd:    "et",
        Template:      "{{range $key, $value := .Attrs }}/v1/kv/services/{{$key}} {{$value}}\n{{end}}",
    }
    templates["ALL"] = append(templates["ALL"], &template)
    parseTemplates(templates)
    
     log.Println("Generate config file %v ", config.ConfigDir(), conf.Filename())

    conf.Save()


    service:= api.Service {
        ID: "idddd",
    	Name: "container A",
    	Port: 8080,
    	IP:    "0.0.0.0",
    	Tags:  make([]string,0),
    	Attrs: map[string]string{
    	    "attr1": "value1",
    	    "attr2": "value2",
    	},
    }
    
    assert.NotNil(t, templates["ALL"])
    assert.NotNil(t, templates["ALL"][0])
    
    query,err := executeTemplates(templates["ALL"][0], &service)
    log.Printf("Esxecuted template: %v ", query)
    assert.Nil(t, err)
    assert.NotNil(t, query)
    assert.Equal(t, "/v1/kv/services/attr1 value1\n/v1/kv/services/attr2 value2\n", query)

}

