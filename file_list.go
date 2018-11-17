package main

import (
	"os"
	"path/filepath"
	"strings"
)

type fileList []os.FileInfo

func (p fileList) filter(fn func(i os.FileInfo) bool) fileList {
	var out fileList
	for _, i := range p {
		if fn(i) {
			out = append(out, i)
		}
	}

	return out
}

func (p fileList) filesMatchingPrefix(prefix string) fileList {
	return p.filter(func(info os.FileInfo) bool {
		return strings.HasPrefix(filepath.Base(info.Name()), prefix)
	})
}

func (p fileList) executables() fileList {
	return p.filter(func(info os.FileInfo) bool {
		return (info.Mode() & 0111) != 0
	})
}

func (p fileList) unique() fileList {
	m := make(map[string]bool)
	return p.filter(func(info os.FileInfo) bool {
		if _, ok := m[info.Name()]; !ok {
			m[info.Name()] = true
			return true
		}

		return false
	})
}
