package main

import (
	"io/ioutil"
	"log"
	"os"
	"net/url"
	"html"
	"html/template"
	"io"
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

func Recurse(path string) Directory {
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
			ret.Directories = append(ret.Directories, Recurse(path + "/" + f.Name()))
		} else if f.Mode().IsRegular() {
			ret.Files = append(ret.Files, Link{
				Link: path + "/" + url.PathEscape(f.Name()),
				Name: html.EscapeString(f.Name()),
			})
		}
	}
	return ret
}

func Render(directories Directory, outf io.Writer) error {
	box := packr.NewBox("./resources")
	index_html := box.String("index.template.html")
	t := template.Must(template.New("index.html").Parse(index_html))
	err := t.Execute(outf, directories)
	if err != nil {
		return err
	}
	return nil
}
	

func main() {
	directories := Recurse(".")

	outf, err := os.Create("index.html")
	if err != nil {
		log.Panic("Could not open index.html for writing. ", err)
	}
	defer outf.Close()	

	err = Render(directories, outf)
	if err != nil {
		log.Panic(err)
	}
}
