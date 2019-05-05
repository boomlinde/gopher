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
	".py": '0', ".h": '0', ".md": '0', ".go": '0',
	".fs": '0',
}

var displaytypes = map[rune]string{
	'0': "TXT", '1': "DIR", 's': "SND", 'g': "GIF",
	'I': "PIC", '9': "BIN", '5': "ARC", 'h': "HTM",
}

func getft(fpath string) (rune, error) {
	f, err := os.Lstat(fpath)
	if err != nil {
		return '9', err
	}
	if f.Mode()&os.ModeSymlink != 0 {
		f, err = os.Stat(fpath)
		if err != nil {
			return '9', err
		}
	}
	if f.IsDir() {
		return '1', nil
	}
	extension := strings.ToLower(filepath.Ext(f.Name()))
	t, ok := filetypes[extension]
	if !ok {
		return '9', nil
	}
	return t, nil
}
