package main

import (
	"log"
	"net/http"
	"os"
)

var (
	airtableApiKey = os.Getenv("AIRTABLE_API_KEY")
	airtableBaseId = os.Getenv("AIRTABLE_BASE_ID")
)

func main() {
	if airtableApiKey == "" {
		panic("AIRTABLE_API_KEY must be specified")
	}

	if airtableBaseId == "" {
		panic("AIRTABLE_BASE_ID must be specified")
	}

	client, err := NewRedirectClient(airtableApiKey, airtableBaseId)
	if err != nil {
		panic(err)
	}

	defaultRedirect, err := client.GetDefault()
	if err != nil {
		panic("A 'default' redirect must be specified")
	}
	log.Println("Default redirect " + defaultRedirect.URL)

	http.HandleFunc("/_ah/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("PONG"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, defaultRedirect.URL, 302)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
