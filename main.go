package main

import (
	"fmt"
	"os"
	"strings"
	version "github.com/hashicorp/go-version"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s root_path\n", os.Args[0])
		os.Exit(1)
	}
	root := os.Args[1]
	pkgs, err := scan(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to scan packages: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d packages\n", len(pkgs))
	for _, p := range pkgs {
		fmt.Printf("%s-%s\n", p.Name, p.Version)
	}
	db, err := fetchSecDb()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to download security database: %v", err)
		os.Exit(1)
	}
	insecurePkgs := make(map[string]SecdbPackage)
	for _, secPkg := range db.Packages {
		insecurePkgs[secPkg.Pkg.Name] = secPkg
	}

	foundInsecure := false

	for _, p := range pkgs {
		currentVersion, err := version.NewVersion(p.Version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse version of %s-%s: %s", p.Name, p.Version, err)
			continue
		}
		if val, ok := insecurePkgs[p.Name]; ok {
			var allCves []string
			for _insecure, cves := range val.Pkg.Secfixes {
				secureVersion, err := version.NewVersion(_insecure)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to parse version %s of security database for %s: %s", _insecure, p.Name, err)
					continue
				}
				if currentVersion.LessThan(secureVersion) {
					allCves = append(allCves, cves...)
				}
			}
		    if len(allCves) > 0 {
				fmt.Printf("insecure package: %s (%s)\n", p.Name, strings.Join(allCves, ","))
				foundInsecure = true
			}
		}
	}
	if !foundInsecure {
		fmt.Printf("no known insecure packages found\n")
	}
}
