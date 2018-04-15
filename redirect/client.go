package redirect

import (
	"fmt"

	"github.com/fabioberger/airtable-go"
)

const (
	// DefaultRedirectKey should be fetched if the requested key is not found.
	DefaultRedirectKey = "default"
	redirectTableName  = "Redirects"
)

// Client uses Airtable to implement Redirector.
type Client struct {
	airtableGo *airtable.Client
}

type airtableRecord struct {
	AirtableID string
	Fields     Redirect
}

// NewClient builds an instance of client.
func NewClient(apiKey string, baseID string) (*Client, error) {
	airtableGo, err := airtable.New(apiKey, baseID)
	return &Client{airtableGo: airtableGo}, err
}

// Get implements Redirector.
func (c *Client) Get(key string) (*Redirect, error) {
	var records []airtableRecord

	err := c.airtableGo.ListRecords(
		redirectTableName,
		&records,
		airtable.ListParameters{FilterByFormula: fmt.Sprintf("{Key} = '%s'", key), MaxRecords: 1},
	)

	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("redirect %s not found", key)
	}

	return &records[0].Fields, nil
}
