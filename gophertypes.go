package main

import (
	"fmt"
	"io"
)

type gopherline struct {
	Ftype rune
	Text  string
	Path  string
	Host  string
	Port  int
}

type gopherdir []*gopherline

func (l *gopherline) serialize(w io.Writer) {
	w.Write([]byte(string(l.Ftype)))
	w.Write([]byte(l.Text + "\t"))
	w.Write([]byte(l.Path + "\t"))
	w.Write([]byte(l.Host + "\t"))
	w.Write([]byte(fmt.Sprintf("%d", l.Port)))
	w.Write([]byte("\r\n"))
}

func (d gopherdir) serialize(w io.Writer) {
	for _, l := range d {
		l.serialize(w)
	}
	w.Write([]byte(".\r\n"))
}

var notfound = gopherdir{
	&gopherline{'i', "File not found", "/", "none", 0},
}
