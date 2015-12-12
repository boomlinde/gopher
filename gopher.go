package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var config struct {
	intport int
	extport int
	name    string
	dir     string
}

var filetypes = map[string]rune{
	".txt": '0', ".gif": 'g', ".jpg": 'I', ".jpeg": 'I',
	".png": 'I', ".html": 'h', ".ogg": 's', ".mp3": 's',
	".wav": 's', ".mod": 's', ".it": 's', ".xm": 's',
	".mid": 's', ".vgm": 's', ".s": '0', ".c": '0',
	".py": '0', ".h": '0', ".md": '0',
}

func getft(f os.FileInfo) rune {
	if f.IsDir() {
		return '1'
	}
	extension := strings.ToLower(filepath.Ext(f.Name()))
	t, ok := filetypes[extension]
	if !ok {
		return '9'
	}
	return t
}

type gopherline struct {
	ftype rune
	text  string
	path  string
	host  string
	port  int
}

type gopherdir []*gopherline

func (l *gopherline) serialize(w io.Writer) {
	w.Write([]byte(string(l.ftype)))
	w.Write([]byte(l.text + "\t"))
	w.Write([]byte(l.path + "\t"))
	w.Write([]byte(l.host + "\t"))
	w.Write([]byte(fmt.Sprintf("%d", l.port)))
	w.Write([]byte("\r\n"))
}

var notfound = gopherdir{
	&gopherline{ftype: 'i', text: "File not found", path: "/", host: "none", port: 0},
}

func (d gopherdir) serialize(w io.Writer) {
	for _, l := range d {
		l.serialize(w)
	}
	w.Write([]byte(".\r\n"))
}

func topath(base, gpath string) string {
	gsplit := strings.Split(gpath, "/")
	rel := strings.Join(gsplit, string(os.PathSeparator))
	return path.Clean(path.Join(base, rel))
}

func gettree(base string) (map[string]bool, error) {
	m := make(map[string]bool)
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if !strings.HasPrefix(info.Name(), ".") {
			m[path] = true
		}
		return err
	})
	return m, err
}

func newselector(line string) *gopherline {
	ret := gopherline{'i', "", "/", config.name, 70}
	ret.ftype, line = rune(line[0]), line[1:]
	fields := strings.Split(line, "\t")
	if len(fields) > 0 {
		ret.text = fields[0]
	}
	if len(fields) > 1 {
		ret.path = fields[1]
	}
	if len(fields) > 2 {
		ret.host = fields[2]
	}
	if len(fields) > 3 {
		ret.port, _ = strconv.Atoi(fields[3])
	}
	return &ret
}

func getheader(p string) gopherdir {
	data, err := ioutil.ReadFile(path.Join(p, ".head"))
	l := gopherdir{}
	if err != nil {
		return l
	}
	lines := strings.Split(strings.Replace(string(data), "\r", "", -1), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			l = append(l, newselector(line))
		}
	}
	return l
}

func getdir(p, sel string) gopherdir {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}
	dir := getheader(p)
	for _, f := range files {
		fname := f.Name()
		if strings.HasPrefix(fname, ".") {
			continue
		}
		line := gopherline{
			ftype: getft(f),
			text:  fname,
			path:  sel + "/" + fname,
			host:  config.name,
			port:  config.extport,
		}
		dir = append(dir, &line)
	}
	return dir
}

func handle(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	sel := ""
	for {
		buf, prefix, err := r.ReadLine()
		if err != nil {
			log.Println(err)
		}
		sel += string(buf)
		if !prefix {
			break
		}
	}
	if !strings.HasPrefix(sel, "/") {
		sel = "/" + sel
	}
	if strings.HasSuffix(sel, "/") {
		sel = sel[:len(sel)-1]
	}
	p := topath(config.dir, sel)
	all, err := gettree(config.dir)
	if err != nil {
		notfound.serialize(conn)
		log.Fatal(err)
	}
	if !all[p] {
		notfound.serialize(conn)
		return
	}
	fi, err := os.Stat(p)
	if err != nil {
		notfound.serialize(conn)
		return
	}
	if fi.IsDir() {
		getdir(p, sel).serialize(conn)
	} else {
		f, err := os.Open(p)
		if err != nil {
			notfound.serialize(conn)
			return
		}
		defer f.Close()
		_, err = io.Copy(conn, f)
		if err != nil {
			return
		}
	}
}

func main() {
	hn, err := os.Hostname()
	if err != nil {
		hn = "localhost"
	}
	flag.IntVar(&config.intport, "i", 7000, "The port to listen to")
	flag.IntVar(&config.extport, "e", 70, "The externally visible port")
	flag.StringVar(&config.dir, "d", ".", "The gopher root dir")
	flag.StringVar(&config.name, "n", hn, "The hostname")
	flag.Parse()
	config.dir = filepath.Clean(config.dir)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", config.intport))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}
