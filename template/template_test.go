package template

import (
    "log"
    "strings"
    "testing"
    "text/template"
    "os"
    "github.com/stretchr/testify/assert"
)

func TestNewEmptyClientUrl(t *testing.T) {
    templ , err := NewTemplate("")
    assert.Nil(t, templ)
    assert.NotNil(t, err)
    
}

func TestNewNoTemplate(t *testing.T) {

    templ, err :=NewTemplate("http://test:2375/")
    assert.NotNil(t, err)
    assert.Nil(t, templ)
    //assert.Equal(t, &TemplateRegistry{templates: templates, url: "http://test:2375/"}, templ)
}

func TestNewWithGoodtemplate(t *testing.T) {
    
    // Set Env Variable with template
    os.Setenv("ALL_TMPL_TEST","my template")

    // Run init
    templ, err :=NewTemplate("http://test:2375/")
    os.Unsetenv("ALL_TMPL_TEST")
    // Check
    assert.Nil(t, err)
    assert.NotNil(t, templ)
    
}

func TestGetTemplates(t *testing.T) {
    
        // Set Env Variable with template
    os.Setenv("ALL_TMPL_TEST","my template")

    // Run init
    templates :=getTemplates()
    os.Unsetenv("ALL_TMPL_TEST")

    assert.True(t,len(templates)>0, "Templates can't be empty")
    
    assert.NotNil(t, templates["ALL"])
}

func TestExecuteTemplate(t *testing.T) {
    templates := make(map[string][]*template.Template)
    templates["ALL"]=append(templates["ALL"], template.Must(template.New("etcd template").Parse("my template {{.}}")))
    
    query :=executeTemplates(templates,"ALL", "No data")
    log.Println("Executed template: ", strings.Join(query[:],","))
    assert.NotNil(t, query)
    assert.Equal(t, 1, len(query))
    assert.Equal(t,"my template No data", query[0])
    

}