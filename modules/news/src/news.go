package main

import (
	"fmt"
	"net/http"

	"wfdr/perms"
	"wfdr/pages"
	"wfdr/template"
)

var defaultPS pages.PageServer = &pages.PageServer{Prefix: "/news/"}

func Handler(c http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		defaultPS.SaveHandler(c, r)
		return
	}
	rp := r.URL.Path
	// /news/<article>/edit triggers this
	if rp[len(rp)-len("/edit"):] == "/edit" {
		defaultPS.EditHandler(c, r)
		return
	}
	if len(rp) != len("/news/") {
		defaultPS.ArticleHandler(c, r)
		return
	}

	plist, err := pages.List()
	if err != nil {
		template.Error500(c, r, err)
		return
	}
	
	template.Render(c, r, "News", "main", plist)
	return
}

func main() {
	fmt.Printf("Loading news server...\n")
	template.SetModuleName("news")
	http.HandleFunc("/news/add", AddHandler)
	http.HandleFunc("/news/", Handler)
	http.ListenAndServe(":8170", nil)
}
