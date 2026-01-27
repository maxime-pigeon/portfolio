package main

import (
	"encoding/xml"
	"flag"
	"net/url"
	"os"
	"path/filepath"
	"text/template"
)

const (
	repoRoot    = "../.."
	projectTmpl = "project.tmpl"
	indexTmpl   = "index.tmpl"
)

var (
	baseURL string

	contentDir = filepath.Join(repoRoot, "content")
	tmplsDir   = filepath.Join(repoRoot, "templates")
	outputDir  = filepath.Join(repoRoot, "docs")

	projectsXML = filepath.Join(contentDir, "projects.xml")

	projectTmplPath = filepath.Join(tmplsDir, projectTmpl)
	indexTmplPath   = filepath.Join(tmplsDir, indexTmpl)
)

func main() {
	flag.StringVar(&baseURL, "b", "", "base URL (for GitHub Pages)")
	flag.Parse()

	site := newSite()
	fm := template.FuncMap{"absURL": absURL}
	site.Templates["index"] = template.New(indexTmpl).Funcs(fm)
	site.Templates["project"] = template.New(projectTmpl).Funcs(fm)

	for i, project := range site.Projects {
		var data pageData
		if i > 0 {
			data.Prev = site.Projects[i-1]
		}
		if i < len(site.Projects)-1 {
			data.Next = site.Projects[i+1]
		}
		data.Project = project
		renderProject(data, site.Templates["project"])
	}
	renderIndex(pageData{Site: *site}, site.Templates["index"])
}

type Site struct {
	BaseURL   string
	Templates map[string]*template.Template
	XMLName   xml.Name  `xml:"projects"`
	Projects  []Project `xml:"project"`
}

type Project struct {
	Title  string   `xml:"title"`
	Slug   string   `xml:"slug"`
	Year   int      `xml:"year"`
	Size   string   `xml:"size"`
	Images []string `xml:"images>source"`
	Desc   string   `xml:"description"`
}

type pageData struct {
	Site Site
	Project
	Next, Prev Project
}

func newSite() *Site {
	data, err := os.ReadFile(projectsXML)
	if err != nil {
		panic(err)
	}
	s := Site{BaseURL: baseURL}
	if err := xml.Unmarshal(data, &s); err != nil {
		panic(err)
	}
	s.Templates = map[string]*template.Template{}
	return &s
}

func renderProject(d pageData, t *template.Template) {
	dir := filepath.Join(outputDir, url.PathEscape(d.Slug))
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}
	file, err := os.Create(filepath.Join(dir, "index.html"))
	if err != nil {
		panic(err)
	}
	if _, err := t.ParseFiles(projectTmplPath); err != nil {
		panic(err)
	}
	t.Execute(file, d)
}

func renderIndex(d pageData, t *template.Template) {
	file, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		panic(err)
	}
	if _, err := t.ParseFiles(indexTmplPath); err != nil {
		panic(err)
	}
	t.Execute(file, d)
}

func absURL(path string) string {
	result, err := url.JoinPath(baseURL, path)
	if err != nil {
		return path
	}
	return result
}
