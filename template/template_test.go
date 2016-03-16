package template

import (
    "testing"
    //"text/template"
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
    
    // Check
    assert.Nil(t, err)
    assert.NotNil(t, templ)
}