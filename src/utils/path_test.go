package utils

import (
	"testing"
)

func TestPathExist(t *testing.T) {
	testCase := []struct {
		path  string
		exist bool
	}{
		{path: ".", exist: true},
		{path: "..", exist: true},
		{path: "z:\\asdfaaaaaaaa\asdf asdf", exist: false},
	}
	for _, c := range testCase {
		if PathExist(c.path) != c.exist {
			t.Fatalf("path %s exist should be %v", c.path, c.exist)
		}
	}
	t.Log("success")
}

func TestPathInfo(t *testing.T) {
	t.Log(PathInfo("."))
}
