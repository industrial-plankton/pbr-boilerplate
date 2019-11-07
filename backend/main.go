package main

import (
	"net/http"
	"strconv"
	"time"

	"backend/postgres"
	// "backend/sampleEndpoint"

	// "strings"

	// "github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/common/log"

	"encoding/json"
	"fmt"
	"io/ioutil"

	//"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type config struct {
	AppEnv string `envconfig:"APP_ENV" default:"development"`
}

// type Handler interface {
// 	ServeHTTP(http.ResponseWriter, *http.Request)
// }

type SubAssem struct {
	Parent string
	Child  string
	Qty    int `db:"telcode"`
}

var (
	port    = "3030"
	mysqlDB *sqlx.DB
	env     config
	//httpHandle Handler
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Println(name, " took ", elapsed)
}

func init() {
	envconfig.MustProcess("pbr", &env)
}

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

func refreshUnit(mysqlDB *sqlx.DB) {
	defer timeTrack(time.Now(), "Unit")
	// fetch all places from the db
	var values [][]string
	rows, _ := mysqlDB.Query("SELECT `parts`.`index_Parts`, `unit`.`unit`, `parts`.`IP SKU` FROM `demodb`.`parts` `parts`, `demodb`.`unit` `unit` WHERE `parts`.`units` = `unit`.`index_unit` ORDER BY `parts`.`index_Parts`")

	// iterate over each row
	for rows.Next() {
		var units string
		var index_Parts int
		var SKU string

		// note that city can be NULL, so we use the NullString type
		_ = rows.Scan(&index_Parts, &units, &SKU)
		rowdata := []string{units, SKU, strconv.Itoa(index_Parts)}
		values = append(values, rowdata)
	}
	sliceofslices := values[:]
	fmt.Println(sliceofslices, "\n")
}

func refreshMaster2(mysqlDB *sqlx.DB) [][]interface{} {
	defer timeTrack(time.Now(), "master2")
	// fetch all places from the db
	var values [][]interface{}
	rows, _ := mysqlDB.Queryx("SELECT `parts`.`IP SKU`, `parts`.`Technical Description`, `parts`.`Customer Description`, `vendors`.`name`, `parts`.`Main Supplier PN` FROM `demodb`.`parts` `parts`, `demodb`.`vendors` `vendors` WHERE `parts`.`Supplier (Main)` = `vendors`.`index_Ven` ORDER BY `parts`.`index_Parts`")
	// iterate over each row
	for rows.Next() {
		var tdisc string
		var cdisc string
		var SKU string
		var mVen string
		var MPN string
		// note that city can be NULL, so we use the NullString type
		_ = rows.Scan(&SKU, &tdisc, &cdisc, &mVen, &MPN)
		rowdata := []interface{}{SKU, tdisc, cdisc, mVen, MPN}
		values = append(values, rowdata)
	}

	var values2 [][]interface{}
	rows, _ = mysqlDB.Queryx("SELECT `parts`.`IP SKU`, `vendors`.`name`, `parts`.`Secondary supplier PN`, `parts`.`Extra Info`, `Order Type`.`Order Type`, `unit`.`unit` FROM `demodb`.`parts` `parts`, `demodb`.`vendors` `vendors`, `demodb`.`Order Type` `Order Type`, `demodb`.`unit` `unit` WHERE `parts`.`Supplier (Secondary)` = `vendors`.`index_Ven` AND `parts`.`Order Type` = `Order Type`.`index_OType` AND `parts`.`units` = `unit`.`index_unit` ORDER BY `parts`.`index_Parts`")

	// iterate over each row
	for rows.Next() {
		var SKU string
		var sVen string
		var SPN string
		var ExtraInf string
		var OT string
		var unit string

		// note that city can be NULL, so we use the NullString type
		_ = rows.Scan(&SKU, &sVen, &SPN, &ExtraInf, &OT, &unit)
		rowdata := []interface{}{sVen, SPN, ExtraInf, OT, unit}
		values2 = append(values2, rowdata)
	}

	for index, element := range values2 {
		// index is the index where we are
		// element is the element from someSlice for where we are
		values[index] = append(values[index], element...)
	}

	return values
}

func writeToSpreadSheet(SQLData [][]interface{}, srv *sheets.Service) {
	defer timeTrack(time.Now(), "Write")
	ctx := context.Background()
	spreadsheetId := "1Hi0PrHe53q4JhNetcJ_y3WrDIJ9qocVEd4irMunxVyE"

	// How the input data should be interpreted.
	valueInputOption := "USER_ENTERED"
	rangeData := "'Master Part List'!A2:Y1400"
	data := SQLData

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: valueInputOption,

		// TODO: Add desired fields of the request body.
	}
	// The new values to apply to the spreadsheet.
	rb.Data = append(rb.Data, &sheets.ValueRange{
		Range:  rangeData,
		Values: data,
	})

	resp, err := srv.Spreadsheets.Values.BatchUpdate(spreadsheetId, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Change code below to process the `resp` object:
	fmt.Printf("%#v\n", resp)

}

func tokenFromString(Auth http.Header) *oauth2.Token {
	// Reads the Authorization JSON sent with the url request and parses it into a oauth2 Token that can be used
	AuthT := &oauth2.Token{}
	AuthT.AccessToken = Auth.Get("AccessToken")
	expiry, err := strconv.ParseInt(Auth.Get("Expiry"), 10, 64)
	if err != nil {
		panic(err)
	}
	AuthT.Expiry = time.Unix(expiry, 0)
	AuthT.RefreshToken = Auth.Get("RefreshToken")
	AuthT.TokenType = "Bearer"
	// decoder := json.NewDecoder(strings.NewReader(Auth))
	// err := decoder.Decode(&AuthT)
	// if err != nil {
	// 	panic(err)
	// }

	return AuthT
}

func getSheetsService(Auth http.Header) *sheets.Service {
	ctx := context.Background()

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// // If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	// client := config.Client(context.Background(), tok)

	// srv, err := sheets.New(client)
	srv, err := sheets.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, tokenFromString(Auth))))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	return srv
}

func test(w http.ResponseWriter, r *http.Request) {
	// Auth := tokenFromString(r.Header.Get("Authorization"))
	Auth := getSheetsService(r.Header)
	values := refreshMaster2(mysqlDB)
	writeToSpreadSheet(values, Auth)
}

func main() {
	// Setup our database connection
	mysqlDB = mysql.NewConnection()

	//	refreshSubAss(mysqlDB)
	//refreshUnit(mysqlDB)
	//values := refreshMaster2(mysqlDB)
	//writeToSpreadSheet(values)
	// Setup our router
	//main := mux.NewRouter()
	// apiSubrouterPath := "/api"
	// routerAPI := main.PathPrefix(apiSubrouterPath).Subrouter()
	// routerV1 := routerAPI.PathPrefix("/v1").Subrouter()

	// Load our endpoints
	//sampleEndpoint.Load(routerV1, mysqlDB)

	log.Info("The server is starting, and it will be listening on port " + port)

	//server := &http.Server{Addr: ":" + port, Handler: routerAPI}
	// server := &http.Server{Addr: ":" + port, Handler: httpHandle}

	// // Prevents memory leak
	// server.SetKeepAlivesEnabled(false)
	// server.ListenAndServe()
	// // HTTP Rest server
	// log.Fatal(
	// 	// Serve on the specified port
	// 	server.ListenAndServe(),
	// )
	http.HandleFunc("/", test)
	if err := http.ListenAndServe(":3030", nil); err != nil {
		panic(err)
	}
}
