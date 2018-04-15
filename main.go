package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/zqureshi/go/redirect"
)

var (
	airtableAPIKey string
	airtableBaseID string
)

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "PONG")
}

func init() {
	airtableAPIKey = os.Getenv("AIRTABLE_API_KEY")
	if airtableAPIKey == "" {
		panic("AIRTABLE_API_KEY must be specified")
	}

	airtableBaseID = os.Getenv("AIRTABLE_BASE_ID")
	if airtableBaseID == "" {
		panic("AIRTABLE_BASE_ID must be specified")
	}
}

func main() {
	c, err := redirect.NewClient(airtableAPIKey, airtableBaseID)
	if err != nil {
		panic(err)
	}

	client, err := redirect.NewCachingClient(c)
	if err != nil {
		panic(err)
	}

	defaultRedirect, err := client.Get(redirect.DefaultRedirectKey)
	if err != nil {
		panic("A 'default' redirect must be specified")
	}
	log.Println("Default redirect " + defaultRedirect.URL)

	http.HandleFunc("/_ah/health", healthHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if r, err := client.Get(req.URL.Path[1:]); req.URL.Path != "/" && err == nil {
			http.Redirect(w, req, r.URL, 302)
		} else {
			d, _ := client.Get(redirect.DefaultRedirectKey)
			http.Redirect(w, req, d.URL, 302)
		}
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}
