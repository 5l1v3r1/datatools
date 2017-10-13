
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

