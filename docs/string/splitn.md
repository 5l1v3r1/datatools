
# string splitn DELIMITER N [JSON_ARRAY]

This command will return split a string at the delimiter
and return it as a JSON arrya.

## Typical command line

```shell
    string splitn '|' 'one|two|three' 2
```

Would return `["one","two|three"]`

## Piping content

NOTE: To read content from standard input we use the `-i -` option.

```shell
    echo 'one|two|three' | string -i - splitn '|' 2
```

Would return `["one","two|three"]`
