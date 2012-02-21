package main

import (
	"fmt"
	"net/http"
	"net/url"

	"wfdr/session"
	"github.com/fduraffourg/go-openid"
)

func Handler(c http.ResponseWriter, r *http.Request) {
	host := r.Host
	s := session.Get(c, r)
	
	continueURL := r.URL.Query().Get("continue-url")
	if continueURL == "" {
		continueURL = "/"
	}
	s.Set("openid-continue-url", continueURL)
	
	// WTF?
	baseUrl := "https://www.google.com/accounts/o8/ud"
	var urlParams = map[string]string{
		"openid.ns":                "http://specs.openid.net/auth/2.0",
		"openid.claimed_id":        "http://specs.openid.net/auth/2.0/identifier_select",
		"openid.identity":          "http://specs.openid.net/auth/2.0/identifier_select",
		"openid.return_to":         "http://" + host + "/openid/auth",
		"openid.realm":             "http://" + host,
		"openid.mode":              "checkid_setup",
		"openid.ns.ui":             "http://specs.openid.net/extensions/ui/1.0",
		"openid.ns.ext1":           "http://openid.net/srv/ax/1.0",
		"openid.ext1.mode":         "fetch_request",
		"openid.ext1.type.email":   "http://axschema.org/contact/email",
		"openid.ext1.type.first":   "http://axschema.org/namePerson/first",
		"openid.ext1.type.last":    "http://axschema.org/namePerson/last",
		"openid.ext1.type.country": "http://axschema.org/contact/country/home",
		"openid.ext1.type.lang":    "http://axschema.org/pref/language",
		"openid.ext1.required":     "email,first,last,country,lang",
		"openid.ns.oauth":          "http://specs.openid.net/extensions/oauth/1.0",
		"openid.oauth.consumer":    host,
		"openid.oauth.scope":       "http://picasaweb.google.com/data/",
	}

	queryURL := "?"
	for name, value := range urlParams {
		queryURL += url.QueryEscape(name) + "=" + url.QueryEscape(value) + "&"
	}
	queryURL = queryURL[0 : len(queryURL)-1]

	http.Redirect(c, r, baseUrl+queryURL, 307)
}

func AuthHandler(c http.ResponseWriter, r *http.Request) {
	grant, _, err := openid.VerifyValues(r.URL.Query())
	if err != nil {
		emsg := fmt.Sprintln("Error in openid auth handler:", err)
		fmt.Println(emsg)
		fmt.Fprintln(c, emsg)
		return
	}
	if !grant {
		fmt.Println("Permission denied!")
		fmt.Fprintln(c, "Access denied by user or internal error.")
		return
	}
	s := session.Get(c, r)
	fmt.Println("Permission granted!")
	
 	wantedValues := []string{"value.email", "value.first", "value.last", "value.country", "value.lang"}
 	qvalues := r.URL.Query()
	for _, wantedValue := range wantedValues {
		value, _ := url.QueryUnescape(qvalues.Get("openid.ext1."+wantedValue))
		s.Set("openid-"+wantedValue[len("value."):], value)
	}
	id, _ := url.QueryUnescape(qvalues.Get("openid.ext1.value.email"))
	s.Set("openid-email", id)
	fmt.Println("OpenID ID:", id)
	
	continueURL := s.Get("openid-continue-url")
	if continueURL == "" {
		continueURL = "/"
	}
	fmt.Println(c, r, continueURL)
	http.Redirect(c, r, continueURL, 307)
	fmt.Fprintln(c, "Authenticated as", id)
	return
}

func WhoamiHandler(c http.ResponseWriter, r *http.Request) {
	s := session.Get(c, r)
	fmt.Fprintln(c, "Authenticated as:", s.Get("openid-email"))
}

func main() {
	fmt.Printf("Loading openid server...\n")
	http.HandleFunc("/openid", Handler)
	http.HandleFunc("/openid/auth", AuthHandler)
	http.HandleFunc("/openid/whoami", WhoamiHandler)
	http.ListenAndServe(":8160", nil)
}
