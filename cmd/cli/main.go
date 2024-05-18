package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"dcashman.net/coctaleague/pkg/bid"
	"dcashman.net/coctaleague/pkg/models"
	"dcashman.net/coctaleague/pkg/models/googlesheets"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {

	// Basic paremeters for the given season
	var (
		maxRuntime int
		numMembers int
		pollFreq   int
		username   string
		sheetRange string
		sheetTitle string
	)

	flag.IntVar(&numMembers, "numMembers", 14, "Number of members in the league")
	flag.IntVar(&maxRuntime, "maxRuntime", 30, "How long to run this program before we stop polling the draft server and making bids.")
	flag.IntVar(&pollFreq, "pollFreq", 30, "How often, in seconds, to poll the draft server and check to see if we need to make a bid")
	flag.StringVar(&username, "username", "Dan", "User for whom to place a bid")
	flag.StringVar(&sheetRange, "range", "A1:DX103", "Range of cells in the spreadsheet, e.g. A2:DX103")
	flag.StringVar(&sheetTitle, "sheetTitle", "Copy of 2023 Draft!", "The sheet to target, e.g. 2023 Draft")

	flag.Parse()

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	// See scopes at https://developers.google.com/identity/protocols/oauth2/scopes
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// TODO: Make this configurable dev vs. prod, right now only dev.
	spreadsheetId := "1hKj3yQduXosNy4Pn9XgzLlqPLArAxKowxip3zKD3UGE"
	readRange := fmt.Sprintf("%s!%s", sheetTitle, sheetRange)

	var draftDb models.DraftStore
	draftDb = googlesheets.NewGoogleSheetsDb(readRange, spreadsheetId, srv, sheetTitle)

	// We've gotten our sheet db, so let's start drafting!  We'll poll the draft until we reach a
	// maximum timeout
	start := time.Now()
	maxDuration := time.Duration(maxRuntime) * time.Second
	for {
		if time.Since(start) > maxDuration {
			// We're done here
			log.Printf("Time's up, shutting down draft polling")
			return
		}

		snapshot, err := draftDb.ParseDraft()
		if err != nil {
			log.Fatalf("Unable to parse draft from Google sheet: %v", err)
		}

		team := snapshot.TeamFromName(username)

		// Check to see if any "basic" bids need to be cast before calculating. These bids are ones
		// which we will always want to opportunistically make, such as making sure we have the best
		// possible player already selected for any of the positions for which we only want to pay one
		// point.
		for _, bid := range bid.RecommendBids(snapshot, team, bid.Value) {
			err := draftDb.PlaceBid(bid)
			if err != nil {
				log.Fatalf("Unable to place Bid: %v. Error: %v", bid, err)
			}
		}

		// Wait for 30 seconds before polling again
		time.Sleep(time.Duration(pollFreq) * time.Second)
	}
}
