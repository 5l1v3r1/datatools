package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/dataset"
)

var (
	// usage
	usage = `USAGE: %s [OPTIONS] JSON_SECRET_FILE SPREADSHEET_ID READ_RANGE`

	// description
	description = `

SYNOPSIS

%s get credentials and read GSheet data and display

`

	// example
	example = `

EXAMPLE

Get GSheet data using a JSON secrets file, the spreadsheet id and 
range.

    %s "etc/client_secret.json" "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" "A2:E"

`

	// Standard Options
	showHelp     bool
	showLicense  bool
	showVersion  bool
	showExamples bool
	inputFName   string
	outputFName  string
	quiet        bool
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	cli.ExitOnError(os.Stderr, err, quiet)
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

	var (
		code string
		err  error
	)
	_, err = fmt.Scan(&code)
	cli.ExitOnError(os.Stderr, err, quiet)

	tok, err := config.Exchange(oauth2.NoContext, code)
	cli.ExitOnError(os.Stderr, err, quiet)
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
	cli.ExitOnError(os.Stderr, err, quiet)
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func init() {
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&showExamples, "example", false, "display example(s)")
	flag.BoolVar(&quiet, "quiet", false, "suppress error messages")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

	cfg := cli.New(appName, appName, dataset.Version)
	cfg.LicenseText = fmt.Sprintf(dataset.License, appName, dataset.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName)
	cfg.OptionText = "OPTIONS\n\n"
	cfg.ExampleText = fmt.Sprintf(example, appName)

	if showHelp == true {
		fmt.Println(cfg.Usage())
		os.Exit(0)
	}

	if showExamples == true {
		fmt.Println(cfg.ExampleText)
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

	fmt.Printf("DEBUG args: %+v\n", args)

	in, err := cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(os.Stderr, err, quiet)
	defer cli.CloseFile(inputFName, in)

	out, err := cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(os.Stderr, err, quiet)
	defer cli.CloseFile(outputFName, out)

	if len(args) != 3 {
		cli.ExitOnError(os.Stderr, fmt.Errorf("expected three parameters, JSON_SECRET_FILE, Sheet ID and Read Range"), quiet)
	}

	jsonSecret := strings.TrimSpace(args[0]) // e.g. etc/client_secret.json

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetId := strings.TrimSpace(args[1]) // e.g. "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"

	readRange := "Class Data!" + strings.TrimSpace(args[2]) // e.g. "Class Data!A2:E"

	ctx := context.Background()

	b, err := ioutil.ReadFile(jsonSecret)
	cli.ExitOnError(os.Stderr, err, quiet)

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	cli.ExitOnError(os.Stderr, err, quiet)
	client := getClient(ctx, config)

	srv, err := sheets.New(client)
	cli.ExitOnError(os.Stderr, err, quiet)

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	cli.ExitOnError(os.Stderr, err, quiet)

	if len(resp.Values) == 0 {
		cli.ExitOnError(os.Stderr, fmt.Errorf("No data found."), quiet)
	}

	var cell string

	w := csv.NewWriter(out)
	for _, values := range resp.Values {
		row := []string{}
		for _, value := range values {
			switch value.(type) {
			case int:
				cell = fmt.Sprintf("%d", value)
			case int64:
				cell = fmt.Sprintf("%d", value)
			case float64:
				cell = fmt.Sprintf("%f", value)
			default:
				cell = fmt.Sprintf("%s", value)
			}
			row = append(row, cell)
		}
		err = w.Write(row)
	}
	cli.ExitOnError(os.Stderr, err, quiet)
	w.Flush()
	err = w.Error()
	cli.ExitOnError(os.Stderr, err, quiet)
}
