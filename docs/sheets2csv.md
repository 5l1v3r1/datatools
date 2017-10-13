
# USAGE

```
    sheets2csv CLIENT_SECRET_JSON SHEETS_ID SHEET_NAME [CELL_RANGE]
```

## Description

sheets2csv accesses a specific sheet id and sheet name in Google Sheets and writes the output as
a CSV file.

It can work from environment settings for the following variables or from options on the commmand

+ GOOGLE_CLIENT_SECRET_JSON (optional) holds the path/filanme to a client secret json file downloaded from Google's developer dashboard
+ GOOGLE_SHEET_ID (optional) holds the hash id for the Google Sheets you want to access
+ GOOGLE_SHEET_NAME (option) holds the specific sheet name in spreadsheet identified by the GOOGLE_SHEET_ID


	-cell-range	set the cell range (e.g. A1:ZZ)
	-client-secret	set path to client secret json
	-example	display example(s)
	-h	display help
	-help	display help
	-l	display license
	-license	display license
	-o	set the filename for output
	-output	set the filename for output
	-sheet-id	set the Google Sheet ID
	-sheet-name	set the Google Sheet Name
	-v	display version
	-version	display version


## Example

In this example we use a Google sheet id is _XOERWASDEWRWEREWRWEFVE23492wd_,
our sheet name is "Untitle 1", we want get cells ranging from "A1" through 
column "G". and out output CSV file will be called
and *data.csv*. We are passing in the client secret json file via the 
environment.

```shell
    export GOOGLE_CLIENT_SECRET_JSON="etc/client_secret.json"
    sheets2csv \
        -cell-range="A1:G" \
        -sheet-name="Untitle 1" \
        -sheet-id=XOERWASDEWRWEREWRWEFVE23492wd \
        -output=data.csv
```

If you want the results to go to *standard out* and specify the 
client secret sheet id, sheet name and cell range in order.

```shell
    export GOOGLE_CLIENT_SECRET_JSON="etc/client_secret.json"
    sheets2csv XOERWASDEWRWEREWRWEFVE23492wd "Untitle 1" "A1:G"
```

sheets2csv v0.0.15
