package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type tplRow struct {
	Link string
	Type string
	Text string
}

func renderHttpMenu(w http.ResponseWriter, tpl *template.Template, d gopherdir) error {
	out := make([]tplRow, len(d))

	for i, x := range d {
		tr := tplRow{Text: x.Text, Type: displaytypes[x.Ftype]}
		if x.Ftype == 'i' {
			out[i] = tr
			continue
		}

		if strings.HasPrefix(x.Path, "URL:") {
			tr.Link = x.Path[4:]
		} else if x.Host == config.name && x.Port == config.extport {
			tr.Link = x.Path
		} else {
			tr.Link = fmt.Sprintf("gopher://%s:%s/%c%s", x.Host, x.Port, x.Ftype, x.Path)
		}

		out[i] = tr
	}

	return tpl.Execute(w, struct {
		Title string
		Lines []tplRow
	}{config.name, out})
}

func serveHttp(addr string) {
	tpldata, err := ioutil.ReadFile(filepath.Join(config.dir, ".template"))
	if err == nil {
		tpltext = string(tpldata)
	}
	tpl, err := template.New("gophermenu").Parse(tpltext)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var p string

		if !contains(config.dir, r.URL.Path) {
			http.Error(w, "File not found", 404)
			return
		}

		p = topath(config.dir, r.URL.Path)

		if isdir(p) {
			if err := renderHttpMenu(w, tpl, getdir(p, r.URL.Path)); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		} else {
			f, err := os.Open(p)
			if err != nil {
				http.Error(w, "File not found", 404)
				return
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				http.Error(w, "File not found", 404)
				return
			}
			w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

			_, err = io.Copy(w, f)
			if err != nil {
				return
			}
		}
	})

	s := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 & time.Second,
		WriteTimeout:   10 & time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go log.Fatal(s.ListenAndServe())
}