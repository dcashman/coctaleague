package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

type ShotclockData struct {
	//lastUpdate time.Time
	duration    time.Duration
	hotseatTeam string
	history     map[string]time.Duration
}

func initializeShotClockData(f string) *ShotclockData {
	return &ShotclockData{
		duration:    0,
		hotseatTeam: "",
		history:     make(map[string]time.Duration),
	}
}

func main() {

	// Basic paremeters for the given season
	var (
		maxRuntime int
		numMembers int
		pollFreq   int
		prod       bool
		username   string
		sheetRange string
		sheetTitle string
		scFile     string
	)

	flag.IntVar(&numMembers, "numMembers", 14, "Number of members in the league")
	flag.IntVar(&maxRuntime, "maxRuntime", 30, "How long to run this program before we stop polling the draft server and making bids.")
	flag.IntVar(&pollFreq, "pollFreq", 30, "How often, in seconds, to poll the draft server and check to see if we need to make a bid")
	flag.StringVar(&username, "username", "Dan", "User for whom to place a bid")
	flag.StringVar(&sheetRange, "range", "DX111", "Second value for range of cells in the spreadsheet, e.g. A1:DX103 should provide DX103. Program starts at A1 by default")
	flag.StringVar(&sheetTitle, "sheetTitle", "2024 Draft", "The sheet to target, e.g. 2023 Draft")
	flag.BoolVar(&prod, "prod", false, "Whether or not to use the real sheet")
	flag.StringVar(&scFile, "scFile", "", "File to record shot-clock time information")

	flag.Parse()

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	// See scopes at https://developers.google.com/identity/protocols/oauth2/scopes
	// Remove readonly for write access.
	//config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// TODO: Make this configurable dev vs. prod, right now only dev.
	spreadsheetId := "1sOHHqvsp4QWZOErmRz0k6MJDxU6fPP_SpUEhWF6BQg8"
	if prod {
		spreadsheetId = "1bzgEDvbHuntqp6FdJiMMg5rmjQ2b5N6pi0BjDy5R8vE"
	}

	sheetRange = fmt.Sprintf("%s:%s", "A1", sheetRange)
	var draftDb models.DraftStore
	draftDb = googlesheets.NewGoogleSheetsDb(sheetRange, spreadsheetId, srv, sheetTitle)

	// Deal with our timekeeping job
	// 1. open the file if one exists to grab the existing values: we need to know the last person who was 'it' to
	// detect if there has been any change.
	// Then we want to write out the value to a file.
	var scData *ShotclockData
	if scFile != "" {
		scData = initializeShotClockData(scFile)
		if err != nil {
			log.Fatalf("Unable to initialize shot clock data: %v", err)
		}
	}

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

		snapshot, err := draftDb.ParseDraft(numMembers)
		if err != nil {
			log.Fatalf("Unable to parse draft from Google sheet: %v", err)
		}

		team, err := models.TeamFromName(snapshot, username)
		if err != nil {
			log.Fatalf("No such team with username: %v", err)
		}

		// TODO: Get from cmdline params.
		bidStrategy := bid.Strategy{Style: bid.Value, Value: bid.Predicted, Preemptive: bid.TwoPointMin}
		for _, bid := range bid.RecommendBids(snapshot, team, bidStrategy) {
			err := draftDb.PlaceBid(bid)
			if err != nil {
				log.Fatalf("Unable to place Bid: %v. Error: %v", bid, err)
			}
		}

		if scData != nil {
			hs := snapshot.Hotseat()
			if strings.Contains(hs, "& other") {
				hs = "Multiple"
			}
			history := snapshot.Times()
			if hs != scData.hotseatTeam {
				scData.hotseatTeam = hs
				scData.duration = 0
			} else {
				scData.duration += time.Second * time.Duration(pollFreq)
				d := history[hs]
				history[hs] = d + time.Second*time.Duration(pollFreq)
			}
			draftDb.WriteShotclock(scData.duration, history[hs], history)
			fmt.Printf("Total times so far:\n%v\n\n", history)
		}

		// Wait for 30 seconds before polling again
		time.Sleep(time.Duration(pollFreq) * time.Second)
	}
}
