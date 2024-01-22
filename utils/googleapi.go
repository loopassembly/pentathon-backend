package utils

import (
	"context"
	"io/ioutil"
	"log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
    spreadsheetID = "1gIOfb8v4GJu1PnMBjxgefQkelb1hyQycQeQjqwFZUvQ"
    readRange = "Sheet4!A:C"
    credentials = "key.json"
)

var sheetsService *sheets.Service
func initializers() {
	
	creds, err := ioutil.ReadFile("key.json")
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to create JWT config: %v", err)
	}

	client := config.Client(context.Background())
	sheetsService, err = sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Google Sheets service: %v", err)
	}
}