package utils

import (
	"context"
	"io/ioutil"
	"sync"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadsheetID = "1gIOfb8v4GJu1PnMBjxgefQkelb1hyQycQeQjqwFZUvQ"
	readRange     = "Sheet4!A:C"
	credentials   = "key.json"
)

type RequestData struct {
	Values [][]interface{} `json:"data"`
}

var sheetsService *sheets.Service
var once sync.Once

func initialize() error {
	creds, err := ioutil.ReadFile("key.json")
	if err != nil {
		return err
	}

	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		return err
	}

	client := config.Client(context.Background())
	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	sheetsService = srv
	return nil
}


func GetSheetsService() (*sheets.Service, error) {
	var err error
	once.Do(func() {
		err = initialize()
	})
	return sheetsService, err
}
