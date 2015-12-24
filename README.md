# jdict
Rich-features Japanese lookup service

## Features
- Words look up by kanji, kana & meaning

## Installation

```sh
$ go build main.go
./main
```

## Populate Data
Data source could be obtained at http://www.edrdg.org/jmdict/j_jmdict.html

```sh
$ ./main populate <data file path>
# Examples:
$ ./main populate data/sample.dict
```

## Query from command lines

```sh
$ ./main query <key string>
# Examples
# Search for kanji
$ ./main query "案件"
{
  "Key": "案件",
  "Entries": [
    {
      "Kanji": [
        {
          "Key": "案件"
        }
      ],
      "Reading": [
        {
          "Key": "あんけん",
          "Pri": [
            "news1",
            "nf11"
          ]
        }
      ],
      "Meaning": [
        {
          "Gloss": [
            "matter in question",
            "subject",
            "case",
            "item"
          ],
          "Meta": {
            "PartOfSpeech": [
              "n"
            ]
          }
        }
      ]
    }
  ]
}
# Search for kana
$ ./main query あんけん
...
# Search for meaning
$ ./main query subject

```

## Query via api interface

```sh
# Start server
$ go run server.go
```

```sh
# Query
$ curl localhost:3000/query/あんけん
```

## Data sources:
+ [Japanese - English] JMDict project http://www.edrdg.org/jmdict/j_jmdict.html
