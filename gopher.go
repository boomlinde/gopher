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
	"time"
)

const selectorcap = 4096
const timeout = 5 * time.Second

var config struct {
	port int
	name string
	dir  string
	http string
}

func topath(base, sel string) string {
	sel = filepath.FromSlash(sel)
	return filepath.Join(base, sel)
}

func contains(sel string) bool {
	return !strings.Contains(filepath.Clean(sel), "..")
}

func newselector(line string) *gopherline {
	ret := gopherline{'i', "", "/", config.name, config.port}
	ret.Ftype, line = rune(line[0]), line[1:]
	fields := strings.Split(line, "\t")
	if len(fields) > 0 {
		ret.Text = fields[0]
	}
	if len(fields) > 1 {
		ret.Path = fields[1]
	}
	if len(fields) > 2 {
		ret.Host = fields[2]
	}
	if len(fields) > 3 {
		ret.Port, _ = strconv.Atoi(fields[3])
	}
	return &ret
}

func getheader(p string) gopherdir {
	l := gopherdir{}

	data, err := ioutil.ReadFile(path.Join(p, ".head"))
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
			Ftype: getft(f),
			Text:  fname,
			Path:  filepath.ToSlash(filepath.Join("/", sel, fname)),
			Host:  config.name,
			Port:  config.port,
		}
		dir = append(dir, &line)
	}
	return dir
}

func isdir(p string) bool {
	fi, err := os.Stat(p)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func handle(conn io.ReadWriteCloser) {
	defer conn.Close()

	r := bufio.NewReader(io.LimitReader(conn, selectorcap))

	buf, toolarge, err := r.ReadLine()
	if err != nil {
		notfound.serialize(conn)
		return
	}
	if toolarge {
		notfound.serialize(conn)
		return
	}

	sel := string(buf)

	if !contains(sel) {
		notfound.serialize(conn)
		return
	}

	p := topath(config.dir, sel)

	if isdir(p) {
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

	flag.IntVar(&config.port, "p", 70, "The port to listen to")
	flag.StringVar(&config.dir, "d", ".", "The gopher root dir")
	flag.StringVar(&config.name, "n", hn, "The hostname")
	flag.StringVar(&config.http, "w", "", "HTTP server address")
	flag.Parse()

	config.dir, err = filepath.Abs(config.dir)
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", config.port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	if config.http != "" {
		go serveHttp(config.http)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		err = conn.SetReadDeadline(time.Now().Add(timeout))
		if err != nil {
			conn.Close()
			continue
		}
		go handle(conn)
	}
}
