package main

import (
	"fmt"
	"net/http"

	"wfdr/pages"
	"wfdr/template"
)

var defaultPS *pages.PageServer = &pages.PageServer{Prefix: "/news/", PageAlias: "Article", Manager: &pages.FSManager{"data/pages/title", "data/pages/content"}}

func Handler(c http.ResponseWriter, r *http.Request) {
	rp := r.URL.Path
	if r.Method == "POST" {
		defaultPS.Save(c, r)
		return
	}
	if rp == "/news/" {
		defaultPS.List(c, r)
		return
	}
	// /news/<article>/edit triggers this
	if rp[len(rp)-len("/edit"):] == "/edit" {
		defaultPS.Edit(c, r)
		return
	}
	if rp == "/news/add" {
		defaultPS.Add(c, r)
		return
	}
	
	defaultPS.Display(c, r)
	return
}

func main() {
	fmt.Printf("Loading news server...\n")
	template.SetModuleName("news")
	http.HandleFunc("/news/", Handler)
	http.ListenAndServe(":8170", nil)
}
