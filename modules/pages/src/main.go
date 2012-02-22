package main

import (
	"fmt"
	"log"
	"net/http"

	"wfdr/perms"
	"wfdr/pages"
	tmpl "wfdr/template"
)

func SaveHandler(c http.ResponseWriter, r *http.Request) {
	p := perms.GetPerms(r)
	if !p.Write {
		fmt.Fprintln(c, "Access Denied")
		return
	}

	name := r.URL.Path[len("/pages/"):]

	content := r.FormValue("content")
	title := r.FormValue("title")

	err := pages.Save(name, title, []byte(content))
	if err != nil {
		tmpl.Error500(c, r, err)
		return
	}
	http.Redirect(c, r, "/pages/"+name, 301)
}

func EditHandler(c http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/pages/") : len(r.URL.Path)-len("/edit")]
	log.Println("Editing page:", name)
	
	page, err := pages.Load(name)
	if err != nil {
		tmpl.Error500(c, r, err)
		return
	}

	log.Println("Request for pages server. Responding.")
	
	tmpl.Render(c, r, "Editing "+page.Title, "edit", page)
}

func ListHandler(c http.ResponseWriter, r *http.Request) {
	plist := pages.List()
	tmpl.Render(c, r, "Pages list", "list", plist)
	return
}

func Handler(c http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		SaveHandler(c, r)
		return
	}
	rp := r.URL.Path
	if rp[len(rp)-len("/edit"):] == "/edit" {
		EditHandler(c, r)
		return
	}
	if rp == "/pages/" {
		ListHandler(c, r)
		return
	}

	pageData, err := pages.Load(r.URL.Path[len("/pages/"):])
	if err != nil {
		// TODO: 404 Error
		tmpl.Error404(c, r, err)
		return
	}
	log.Println("Request for pages server. Responding.")
	
	tmpl.Render(c, r, pageData.Title, "main", pageData)
}

func main() {
	fmt.Printf("Loading pages server...\n")
	tmpl.SetModuleName("pages")
	
	http.HandleFunc("/pages/", Handler)
	http.ListenAndServe(":8150", nil)
}
