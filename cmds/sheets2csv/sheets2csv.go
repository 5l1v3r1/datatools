//
// sheets2csv is a command line program that reads a Google Sheets's sheet and
// renders out a CSV locally. It is inspired by the quickstart example found on
// developer sites for the Google Sheets API v4.
//
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/datatools"

	// Google Sheets packages
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

var (
	// Standard Options
	showHelp     bool
	showLicense  bool
	showExamples bool
	showVersion  bool
	outputFName  string

	// Application Options
	clientSecretJSON string
	spreadSheetId    string
	sheetName        string
	cellRange        string
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("sheets.googleapis.com-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func sheet2CSV(out *os.File, clientSecretJSON, spreadSheetId, sheetName, cellRange string) error {
	ctx := context.Background()

	b, err := ioutil.ReadFile(clientSecretJSON)
	if err != nil {
		return fmt.Errorf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := sheets.New(client)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Sheets Client %v", err)
	}

	// Prints the columns from sheet described by spreadSheetId
	readRange := fmt.Sprintf("%s!%s", sheetName, cellRange)
	resp, err := srv.Spreadsheets.Values.Get(spreadSheetId, readRange).Do()
	if err != nil {
		return fmt.Errorf("Unable to retrieve data from sheet. %v", err)
	}

	if len(resp.Values) > 0 {
		// NOTE: Writing the values out as CSV using encoding/csv package
		table := [][]string{}
		for _, row := range resp.Values {
			cells := []string{}
			for _, val := range row {
				cell := val.(string)
				cells = append(cells, cell)
			}
			table = append(table, cells)
		}
		// Now that we've copied the values, write the table
		w := csv.NewWriter(out)
		w.WriteAll(table)
		if err := w.Error(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("No data found")
	}
	return nil
}

func init() {
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showExamples, "example", false, "display example(s)")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.StringVar(&outputFName, "o", "", "set the filename for output")
	flag.StringVar(&outputFName, "output", "", "set the filename for output")

	// Application Options
	flag.StringVar(&clientSecretJSON, "client-secret", "", "set path to client secret json")
	flag.StringVar(&spreadSheetId, "sheet-id", "", "set the Google Sheet ID")
	flag.StringVar(&sheetName, "sheet-name", "", "set the Google Sheet Name")
	flag.StringVar(&cellRange, "cell-range", "", "set the cell range (e.g. A1:ZZ)")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

	cfg := cli.New(appName, "GOOGLE", datatools.Version)
	cfg.LicenseText = fmt.Sprintf(datatools.LicenseText, appName, datatools.Version)
	cfg.UsageText = fmt.Sprintf("%s", Help["usage"])
	cfg.DescriptionText = fmt.Sprintf("%s", Help["description"])
	cfg.ExampleText = fmt.Sprintf("%s", Help["examples"])

	if showHelp == true {
		if len(args) > 0 {
			fmt.Println(cfg.Help(args...))
		} else {
			fmt.Println(cfg.Usage())
		}
		os.Exit(0)
	}

	if showExamples == true {
		if len(args) > 0 {
			fmt.Println(cfg.Example(args...))
		} else {
			fmt.Println(cfg.ExampleText)
		}
		os.Exit(0)
	}

	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}

	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}

	out, err := cli.Create(outputFName, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer cli.CloseFile(outputFName, out)

	if len(args) < 1 || len(args) > 4 {
		fmt.Println(cfg.Usage())
		os.Exit(1)
	}

	if clientSecretJSON == "" {
		clientSecretJSON = os.Getenv("GOOGLE_CLIENT_SECRET_JSON")
	}
	if spreadSheetId == "" {
		spreadSheetId = os.Getenv("GOOGLE_SHEET_ID")
		if spreadSheetId == "" {
			spreadSheetId, args = cli.PopArg(args)
		}
	}
	if sheetName == "" {
		sheetName, args = cli.PopArg(args)
	}
	if cellRange == "" {
		if args != nil {
			cellRange, args = cli.PopArg(args)
		} else {
			cellRange = "A1:ZZ"
		}
	}

	if err := sheet2CSV(out, clientSecretJSON, spreadSheetId, sheetName, cellRange); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}
