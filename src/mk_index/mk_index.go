package main

import (
	"io/ioutil"
	"log"
	"os"
	"net/url"
	"html"
	"html/template"
	"strings"
	P "path"
	
	"github.com/gobuffalo/packr"
)

type Directory struct {
	Path        string
	Name        string
	Files       []Link
	Directories []Directory
}

type Link struct {
	Name string
	Link string
}

func clean(p string) string {
	if (strings.HasPrefix(p, "./")) {
		return url.PathEscape(p[2:])
	}
	return url.PathEscape(p)
}

func recurse(path string) Directory {
	ret := Directory{
		path,
		P.Base(path),
		[]Link{},
		[]Directory{},
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Panic(err)
	}

	for _, f := range files {
		if f.Mode().IsDir() {
			ret.Directories = append(ret.Directories, recurse(path + "/" + f.Name()))
		} else if f.Mode().IsRegular() {
			ret.Files = append(ret.Files, Link{
				Link: clean(path + "/" + f.Name()),
				Name: html.EscapeString(f.Name()),
			})
		}
	}
	return ret
}

func main() {
	directories := recurse(".")

	outf, err := os.Create("index.html")
	if err != nil {
		log.Panic("Could not open index.html for writing. ", err)
	}
	defer outf.Close()	

	box := packr.NewBox("./resources")
	
	index_html := box.String("index.template.html")
	t := template.Must(template.New("index.html").Parse(index_html))
	err = t.Execute(outf, directories)
	if err != nil {
		log.Panic("Could not execute template", err)
	}
}
