package template

import (
    "bytes"
    "errors"
    "io/ioutil"
    "log"
    "net/url"
    "net/http"
    "os"
    "strings"
    "text/template"
    api "github.com/docking-tools/register/api"
)

// Usage:
// @TODO

type TemplateRegistry struct {
    templates map[string][]*template.Template
    url string // uri path
}

func NewTemplate(uri string) (api.RegistryAdapter, error) {
    u, err := url.ParseRequestURI(uri)
    
    if err!=nil {
        return nil, err
    }

	// Find all environment variables that contain ETCD_ and turn them into templates
	templates := make(map[string][]*template.Template)
	
	for _, env := range os.Environ() {
	    envArray := strings.SplitN(env, "=", 2)
	    text := envArray[1]
	    key := envArray[0]
	    action := "ALL"
	    
	    // extract specific command template (start with REGISTER_TMPL, DEREGISTER_TMPL, PING_TMPL, UPDATE_TMPL, ALL_TMPL)
	    if strings.Contains(key, "_TMPL_") {
	        action = strings.SplitN(key, "_", 2)[0]
	        log.Println("New template: %v", text)
	        templates[action] = append(templates[action], template.Must(template.New("etcd template").Parse(text)))
	    }
	}
	
	if len(templates)==0 {
	    return nil, errors.New("No template found.")
	}
	return &TemplateRegistry{templates: templates, url: u.String()}, nil
}

func (r *TemplateRegistry) Size() int {
    return len(r.templates)
}

func (r *TemplateRegistry) Ping() error {
    return nil
}

func (r *TemplateRegistry) Register(service *api.Service) error {
    toSet, err := r.executeTemplates("REGISTER", service)
    
    if err != nil {
		return err
	}

	for key, value := range toSet {
	    client := &http.Client{}
	    request, err := http.NewRequest("PUT", r.url+key, strings.NewReader(value))
	    request.ContentLength = int64(len(value))
	    
	    response, err := client.Do(request)
	    if (err != nil) {
	        log.Fatal(err)
	        return nil
	    }
	    defer response.Body.Close()
	    contents, err := ioutil.ReadAll(response.Body)
	    if err != nil {
	        log.Fatal(err)
	        return nil
	    }
	    log.Print(contents)
	}

	return nil    
}

func (r *TemplateRegistry) Deregsiter(service *api.Service) error {
    return nil
}

func (r *TemplateRegistry) Update(service *api.Service) error {
    return nil
}

func (r *TemplateRegistry) executeTemplates(action string, service *api.Service) (map[string]string, error) {
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