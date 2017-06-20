package main

import (
	"bytes"
	"testing"
)

func TestSerializeLine(t *testing.T) {
	line := gopherline{
		Ftype: 'h',
		Text:  "abc",
		Path:  "def",
		Host:  "xyz",
		Port:  70,
	}
	expected := "habc\tdef\txyz\t70\r\n"

	buf := new(bytes.Buffer)

	line.serialize(buf)
	if buf.String() != expected {
		t.Error("gopherline.serialize does not match expected")
	}
}

func TestSerializeDir(t *testing.T) {
	dir := gopherdir{
		&gopherline{'h', "abc", "def", "xyz", 70},
		&gopherline{'i', "123", "456", "789", 80},
	}
	expected := "habc\tdef\txyz\t70\r\ni123\t456\t789\t80\r\n.\r\n"

	buf := new(bytes.Buffer)

	dir.serialize(buf)
	if buf.String() != expected {
		t.Error("gopherdir.serialize does not match expected")
	}
}

func TestTopath(t *testing.T) {
	cases := [][3]string{
		[3]string{"/x/y", "z", "/x/y/z"},
		[3]string{"x/y/", "z/a/b/c", "x/y/z/a/b/c"},
	}

	for _, c := range cases {
		if topath(c[0], c[1]) != c[2] {
			t.Errorf("topath failed: (%s, %s) -> %s", c[0], c[1], c[2])
		}
	}
}

func TestContains(t *testing.T) {
	ok_cases := []string{
		"a",
		"a/../",
		"a/b/../",
		"a/b/../../",
	}

	fail_cases := []string{
		"../a",
		"../",
		"../../",
		"a/b/../../../",
	}

	for _, c := range ok_cases {
		if !contains(c) {
			t.Errorf("contains failed: (%s) -> false", c)
		}
	}
	for _, c := range fail_cases {
		if contains(c) {
			t.Errorf("contains failed: (%s) -> true", c)
		}
	}
}
