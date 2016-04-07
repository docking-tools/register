package config

import (
    "text/template"    
)

type ConfigTemplate struct {
    Name        string                      `json:"name"`
    Event       string                      `json:"event"`
    HttpCmd     string                      `json:"httpCmd"`
    Template       string                      `json:"query"`
    tmpl   *template.Template           // Note: not serialized - for internal use only
}

// Filename returns the name of the configuration file
func (config *ConfigTemplate) SetTmpl(tmpl *template.Template)  {
	config.tmpl = tmpl
}

func (config *ConfigTemplate) Tmpl() *template.Template  {
	return config.tmpl
}