package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Mirror struct {
	dir  string
	port int
}

func (m Mirror) Start() {
	log.Fatal(
		http.ListenAndServe(fmt.Sprint(":", m.port), http.FileServer(http.Dir(m.dir))),
	)
}

func enableMirror(flag bool) {
	if flag {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "."
		}
		mirror := Mirror{
			dir:  cwd,
			port: 8080,
		}
		mirror.Start()
	}
}
