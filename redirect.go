package main

import (
	"fmt"
	"github.com/fabioberger/airtable-go"
)

const redirectTableName = "Redirects"
const defaultRedirectKey = "default"

type RedirectClient struct {
	airtableGo *airtable.Client
}

type Redirect struct {
	Key string
	URL string
}

func NewRedirectClient(apiKey string, baseId string) (*RedirectClient, error) {
	airtableGo, err := airtable.New(apiKey, baseId)
	return &RedirectClient{airtableGo: airtableGo}, err
}

func (client *RedirectClient) Get(key string) (*Redirect, error) {
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

func (client *RedirectClient) GetDefault() (*Redirect, error) {
	return client.Get(defaultRedirectKey)
}
