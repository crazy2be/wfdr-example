package main

import (
	"fmt"
	"net/http"

	"wfdr/pages"
	tmpl "wfdr/template"
)

var defaultManager pages.Manager = &pages.FSManager{"data/pages/title", "data/pages/content"}
var defaultPS *pages.PageServer = &pages.PageServer{Prefix: "/pages/", PageAlias: "Page", Manager: defaultManager}

func Handler(c http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		defaultPS.Save(c, r)
		return
	}
	rp := r.URL.Path
	if rp[len(rp)-len("/edit"):] == "/edit" {
		defaultPS.Edit(c, r)
		return
	}
	if rp == "/pages/" {
		defaultPS.List(c, r)
		return
	}
	defaultPS.Display(c, r)
	return
}

func main() {
	fmt.Printf("Loading pages server...\n")
	tmpl.SetModuleName("pages")
	
	http.HandleFunc("/pages/", Handler)
	http.ListenAndServe(":8150", nil)
}
