package main

import (
	"github.com/fabioberger/airtable-go"
	"log"
	"net/http"
	"os"
)

var (
	airtableApiKey    = os.Getenv("AIRTABLE_API_KEY")
	airtableBaseId    = os.Getenv("AIRTABLE_BASE_ID")
	airtableTableName = "Redirects"
)

type Redirect struct {
	Key string
	URL string
}

type Record struct {
	AirtableID string
	Fields     Redirect
}

func main() {
	if airtableApiKey == "" {
		panic("AIRTABLE_API_KEY must be specified")
	}

	if airtableBaseId == "" {
		panic("AIRTABLE_BASE_ID must be specified")
	}

	client, err := airtable.New(airtableApiKey, airtableBaseId)
	if err != nil {
		panic(err)
	}

	records := []Record{}
	err = client.ListRecords(
		airtableTableName,
		&records,
		airtable.ListParameters{FilterByFormula: "{Key} = 'default'", MaxRecords: 1},
	)

	if err != nil {
		panic(err)
	}

	if len(records) != 1 {
		panic("A 'default' redirect must be specified")
	}

	defaultRedirect := records[0].Fields
	log.Println("Default redirect " + defaultRedirect.URL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, defaultRedirect.URL, 302)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
