package IPSheets

import (
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	// "os"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func tokenFromString(Auth http.Header) *oauth2.Token {
	// Reads the Authorization JSON sent with the url request and parses it into a oauth2 Token that can be used
	AuthT := &oauth2.Token{}
	//Convert Headers token parameters to an actual token
	AuthT.AccessToken = Auth.Get("AccessToken")
	expiry, err := strconv.ParseInt(Auth.Get("Expiry"), 10, 64)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	AuthT.Expiry = time.Unix(expiry, 0)
	AuthT.RefreshToken = Auth.Get("RefreshToken")
	AuthT.TokenType = "Bearer" //Tokens are always Bearer Type
	// decoder := json.NewDecoder(strings.NewReader(Auth))
	// err := decoder.Decode(&AuthT)
	// if err != nil {
	// 	panic(err)
	// }

	return AuthT
}

func GetSheetsService(Auth http.Header) *sheets.Service {
	ctx := context.Background()

	//Read the App Credentials File
	b, err := ioutil.ReadFile("/home/cameron/Go_Server/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	//create the configuration with spreadsheets permissions
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	//create the sheets service, used to read and write to sheets
	srv, err := sheets.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, tokenFromString(Auth))))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return srv
}

/////OBSOLETE Saving token on server
//Replaced with passing token from google sheets

// // Retrieve a token, saves the token, then returns the generated client.
// func getClient(config *oauth2.Config) *http.Client {
// 	// The file token.json stores the user's access and refresh tokens, and is
// 	// created automatically when the authorization flow completes for the first
// 	// time.
// 	tokFile := "token.json"
// 	tok, err := tokenFromFile(tokFile)
// 	if err != nil {
// 		tok = getTokenFromWeb(config)
// 		saveToken(tokFile, tok)
// 	}
// 	return config.Client(context.Background(), tok)
// }

// // Request a token from the web, then returns the retrieved token.
// func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
// 	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
// 	fmt.Printf("Go to the following link in your browser then type the "+
// 		"authorization code: \n%v\n", authURL)

// 	var authCode string
// 	if _, err := fmt.Scan(&authCode); err != nil {
// 		log.Fatalf("Unable to read authorization code: %v", err)
// 	}

// 	tok, err := config.Exchange(context.TODO(), authCode)
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve token from web: %v", err)
// 	}
// 	return tok
// }

// // Retrieves a token from a local file.
// func tokenFromFile(file string) (*oauth2.Token, error) {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	tok := &oauth2.Token{}
// 	err = json.NewDecoder(f).Decode(tok)
// 	return tok, err
// }

// // Saves a token to a file path.
// func saveToken(path string, token *oauth2.Token) {
// 	fmt.Printf("Saving credential file to: %s\n", path)
// 	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		log.Fatalf("Unable to cache oauth token: %v", err)
// 	}
// 	defer f.Close()
// 	json.NewEncoder(f).Encode(token)
// }
