package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/alwindoss/casa"
	"github.com/alwindoss/casa/internal/forms"
	"github.com/justinas/nosurf"
)

func AddDefaultData(r *http.Request, td *TemplateData) *TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

type TemplateData struct {
	CSRFToken       string
	StringSlice     []string
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
	Title           string
	InfoMsg         string
	WarnMsg         string
	ErrMsg          string
}

func (ph pageHandler) renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data *TemplateData) {
	data = AddDefaultData(r, data)
	t, ok := ph.templateCache[tmpl]
	if !ok {
		err := fmt.Errorf("unable to find %s in template cache", tmpl)
		log.Fatal(err)
	}
	buff := new(bytes.Buffer)
	err := t.Execute(buff, data)
	if err != nil {
		log.Fatal(err)
	}
	_, err = buff.WriteTo(w)
	if err != nil {
		log.Printf("error writing template to browser: %v", err)
		return
	}

}

var functions = template.FuncMap{
	// The name "inc" is what the function will be called in the template text.
	"inc": func(i int) int {
		return i + 1
	},
	"marshal": func(v interface{}) template.JS {
		a, _ := json.Marshal(v)
		return template.JS(a)
	},
}

// createTemplateCache creates a template cache as a map
func createTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := fs.Glob(casa.FS, "templates/*.page.tmpl")
	if err != nil {
		log.Printf("error glob: %v", err)
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFS(casa.FS, page)
		if err != nil {
			log.Printf("unable to ParseFiles: %v", err)
			return nil, err
		}
		matches, err := fs.Glob(casa.FS, "templates/*.layout.tmpl")
		if err != nil {
			log.Printf("unable to fs.Glob: %v", err)
			return nil, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseFS(casa.FS, "templates/*.layout.tmpl")
			if err != nil {
				log.Printf("unable to ParseGlob: %v", err)
				return nil, err
			}
		}
		myCache[name] = ts
	}
	log.Printf("Cached the templates")
	return myCache, err
}
