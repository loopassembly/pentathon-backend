package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

func main() {
    
    creds, err := ioutil.ReadFile(credentials)
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

    
	// Read data here
    http.HandleFunc("/read", ReadData)
    http.HandleFunc("/create", CreateData)
    
    // Start the HTTP server.
    port := ":8080"
    fmt.Printf("Server is listening on port %s...\n", port)
    http.ListenAndServe(port, nil)
}


func ReadData(w http.ResponseWriter, r *http.Request) {
    resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(r.Context()).Do()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    data, _ := json.Marshal(resp.Values)
    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func CreateData(w http.ResponseWriter, r *http.Request) {
    // Parse request body to get data to be added.
    type RequestData struct {
        Values [][]interface{} `json:"data"`
    }

    var requestData RequestData

    err := json.NewDecoder(r.Body).Decode(&requestData)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    values := sheets.ValueRange{Values: requestData.Values}
    _, err = sheetsService.Spreadsheets.Values.Append(spreadsheetID, readRange, &values).ValueInputOption("RAW").Do()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}