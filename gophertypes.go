package main

import (
	"fmt"
	"io"
)

type gopherline struct {
	Ftype rune   `json:"ftype"`
	Text  string `json:"text"`
	Path  string `json:"path"`
	Host  string `json:"host"`
	Port  int    `json:"port"`
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

var notfound = gopherdir{
	&gopherline{'i', "File not found", "/", "none", 0},
}

func (d gopherdir) serialize(w io.Writer) {
	for _, l := range d {
		l.serialize(w)
	}
	w.Write([]byte(".\r\n"))
}
