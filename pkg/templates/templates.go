package templates

import (
	"embed"
	"errors"
	"html/template"
	"io"
	"log"
)

//go:embed html/*
var templates embed.FS

//go:embed static
var static embed.FS

//NOTE: I've though about renderPage function that renders a full page.
// And then renderTemplate that can render just a standalone ...component
// Or as alternative just Render function with veriodic ...name argument,
// which can render any combination. Probobly better solution.

type Templates struct {
	templates *template.Template
	Static    embed.FS
	FS        embed.FS
}

func (t *Templates) Render(w io.Writer, data interface{}, patterns ...string) error {

	tmpl, err := template.ParseFS(t.FS, prefixNames(patterns)...)

	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}

func (t *Templates) RenderPage(w io.Writer, data interface{}, name string) error {
	var patterns []string

	patterns = append(patterns, "page.html", "navbar.html", "components.html", name)

	return t.Render(w, data, patterns...)

}

func (t *Templates) RenderComponent(w io.Writer, data interface{}, name string) error {

	tmpls, err := template.ParseFS(t.FS, "html/components.html")

	if err != nil {
		return err
	}

	tmpl := tmpls.Lookup(name)

	if tmpl == nil {
		return errors.New("component not found: " + name)
	}

	return tmpl.Execute(w, data)
}

func prefixNames(args []string) []string {
	for i, p := range args {
		args[i] = "html/" + p
	}

	return args
}

func NewTemplate() *Templates {
	tmpls, err := template.ParseFS(templates, "html/*")

	if err != nil {
		log.Fatal(err)
	}
	return &Templates{
		templates: tmpls,
		FS:        templates,
		Static:    static,
	}
}
