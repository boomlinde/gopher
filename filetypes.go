package main

import (
	"os"
	"path/filepath"
	"strings"
)

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
