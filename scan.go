package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type Package struct {
	// unique ID of this package. this will be created as discovered by the library
	// and used for persistence and hash map indexes
	ID string `json:"id"`
	// the name of the package
	Name string `json:"name"`
	// the version of the package
	Version string `json:"version"`
	// if type is a binary package a source package maybe present which built this binary package.
	// must be a pointer to support recursive type:
	Source *Package `json:"source,omitempty"`
	// the file system path or prefix where this package resides
	PackageDB string `json:"-"`
	// a hint on which repository this package was downloaded from
	RepositoryHint string `json:"-"`
	// Module and stream which this package is part of
	Module string `json:"module,omitempty"`
	// Package architecture
	Arch string `json:"arch,omitempty"`
}

const installedFile = "lib/apk/db/installed"

func scan(root string) ([]*Package, error) {
	fullpath := filepath.Join(root, installedFile)
	content, err := ioutil.ReadFile(fullpath)
	if err != nil {
		fmt.Errorf("cannot read %s: %s", fullpath, err)
	}
	pkgs := []*Package{}
	srcs := make(map[string]*Package)

	// It'd be great if we could just use the textproto package here, but we
	// can't because the database "keys" are case sensitive, unlike MIME
	// headers. So, roll our own entry and header splitting.
	var delim = []byte("\n\n")
	entries := bytes.Split(content, delim)
	for _, entry := range entries {
		if len(entry) == 0 {
			continue
		}
		p := Package{
			PackageDB: installedFile,
		}
		r := bytes.NewBuffer(entry)
		for line, err := r.ReadBytes('\n'); err == nil; line, err = r.ReadBytes('\n') {
			l := string(bytes.TrimSpace(line[2:]))
			switch line[0] {
			case 'P':
				p.Name = l
			case 'V':
				p.Version = l
			case 'c':
				p.RepositoryHint = l
			case 'A':
				p.Arch = l
			case 'o':
				if src, ok := srcs[l]; ok {
					p.Source = src
				} else {
					p.Source = &Package{
						Name: l,
					}
					if p.Version != "" {
						p.Source.Version = p.Version
					}
					srcs[l] = p.Source
				}
			}
		}
		pkgs = append(pkgs, &p)
	}
	return pkgs, nil
}
