package main

import (
    "bytes"
    "net/url"
    "os"
    "strings"
    "text/template"
)

// Usage:
// @TODO

type TemplateRegistry struct {
    templates map[string][]*template.Template
    path string // uri path
}

func New(uri *url.URL) RegistryAdapter {
    urls := make([]string, 0)
	if uri.Host != "" {
		urls = append(urls, "http://"+uri.Host)
	}
	
	// Find all environment variables that contain ETCD_ and turn them into templates
	templates := make(map[string][]*template.Template)
	
	for _, env := range os.Environ() {
	    envArray := strings.SplitN(env, "=", 2)
	    text := envArray[1]
	    key := envArray[0]
	    action := "ALL"
	    
	    // extract specific command template (start with REGISTER_TMPL, DEREGISTER_TMPL, PING_TMPL, UPDATE_TMPL)
	    if strings.Contains(key, "_TMPL_") {
	        action = strings.SplitN(key, "_", 2)[0]
	    }
	    
	   templates[action] = append(templates[action], template.Must(template.New("etcd template").Parse(text)))
	    
	}
	return &TemplateRegistry{}
}

func (r *TemplateRegistry) Ping() error {
    return nil
}

func (r *TemplateRegistry) Register(service *Service) error {
    return nil
}

func (r *TemplateRegistry) Deregsiter(service *Service) error {
    return nil
}

func (r *TemplateRegistry) Update(service *Service) error {
    return nil
}

func (r *TemplateRegistry) executeTemplates(action string, service *Service) (map[string]string, error) {
    tmpls := []*template.Template {}
    
    buf := &bytes.Buffer {}
    tmpls = append(tmpls, r.templates[action]...)
    tmpls = append(tmpls, r.templates["ALL"]...) 

    results := make(map[string]string, len(tmpls))
    
    for _,t := range tmpls {
        // Execute the template with the service as the data item
        buf.Reset()
        err := t.Execute(buf, service)
        if err != nil {
            return nil, err
        }
        
        // The templates needs to return "<keyPath> <data>". The key must not 
        // contain any spaces, so we use the first space as the split between the two.
        // If nothing is resturned, then that says not tu use that template
        pair := strings.SplitN(buf.String(), " ", 2)
        if 2== len(pair) {
            key := strings.TrimSpace(pair[0])
            value := strings.TrimSpace(pair[1])
            if len(key) > 0 {
                results[key] = value
            }
        }
    }
        
    return results, nil
}