package redirect

import (
	"fmt"

	"github.com/fabioberger/airtable-go"
)

const redirectTableName = "Redirects"
const defaultRedirectKey = "default"

type Client struct {
	airtableGo *airtable.Client
}

type Redirect struct {
	Key string
	URL string
}

func New(apiKey string, baseId string) (*Client, error) {
	airtableGo, err := airtable.New(apiKey, baseId)
	return &Client{airtableGo: airtableGo}, err
}

func (client *Client) Get(key string) (*Redirect, error) {
	records := []struct {
		AirtableID string
		Fields     Redirect
	}{}

	err := client.airtableGo.ListRecords(
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

func (client *Client) GetDefault() (*Redirect, error) {
	return client.Get(defaultRedirectKey)
}
