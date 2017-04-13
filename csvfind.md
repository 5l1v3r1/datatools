
# USAGE

    csvfind [OPTIONS] TEXT_TO_MATCH

## SYNOPSIS

csvfind processes a CSV file as input returning rows that contain the column
with matched text. Supports exact match as well as some Levenshtein matching.

## OPTIONS

```
	-append-edit-distance	append column with edit distance found (useful for tuning levenshtein)
	-case-sensitive	perform a case sensitive match (default is false)
	-col	column to search for match in the CSV file
	-delete-cost	set the delete cost to use for levenshtein matching
	-h	display help
	-help	display help
	-i	input filename
	-input	input filename
	-insert-cost	set the insert cost to use for levenshtein matching
	-l	display license
	-levenshtein	use levenshtein matching
	-license	display license
	-max-edit-distance	set the edit distance thresh hold for match, default 0
	-o	output filename
	-output	output filename
	-skip-header-row	skip the header row
	-stop-words	use the colon delimited list of stop words
	-substitute-cost	set the substitution cost to use for levenshtein matching
	-v	display version
	-version	display version
```

## EXAMPLES

Find the rows where the third column matches "The Red Book of Westmarch" exactly

```shell
    csvfind -i books.csv -col=2 "The Red Book of Westmarch"
```

Find the rows where the third column matches approximately "The Red Book of Westmarch"

```shell
    csvfind -i books.csv -col=2 -levenshtein \
       -insert-cost=1 -delete-cost=1 -substitute-cost=3 \
       -max-edit-distance=50 -append-edit-distance \
       "The Red Book of Westmarch"
```

In this example all records from the demo books.csv file would be returned with their
distance number as the final column.
