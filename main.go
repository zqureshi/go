package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/zqureshi/go/redirect"
)

var (
	airtableApiKey = os.Getenv("AIRTABLE_API_KEY")
	airtableBaseId = os.Getenv("AIRTABLE_BASE_ID")
)

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "PONG")
}

func main() {
	if airtableApiKey == "" {
		panic("AIRTABLE_API_KEY must be specified")
	}
	if airtableBaseId == "" {
		panic("AIRTABLE_BASE_ID must be specified")
	}

	client, err := redirect.NewCaching(airtableApiKey, airtableBaseId)
	if err != nil {
		panic(err)
	}

	defaultRedirect, err := client.GetDefault()
	if err != nil {
		panic("A 'default' redirect must be specified")
	}
	log.Println("Default redirect " + defaultRedirect.URL)

	http.HandleFunc("/_ah/health", healthHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if redirect, err := client.Get(r.URL.Path[1:]); r.URL.Path != "/" && err == nil {
			http.Redirect(w, r, redirect.URL, 302)
		} else {
			defaultRedirect, err = client.GetDefault()
			http.Redirect(w, r, defaultRedirect.URL, 302)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
