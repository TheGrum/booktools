#  Booktools  #
===============

[![GoDoc](https://godoc.org/github.com/TheGrum/booktools?status.svg)](https://godoc.org/github.com/TheGrum/booktools)

Install:
```
go get github.com/TheGrum/booktools
cd github.com/TheGrum/booktools/booktools
go build
```

=====

Booktools is a simple tool to read a raw text file containing a manuscript 
and parse its structure, to then provide information about it.

### Usage

```
Usage:
  booktools process [command]

Available Commands:
  chapter              Displays the contents of a chapter
  chapterCharacters    Lists the characters in each chapter
  characterFrequencies Lists the characters and the frequency with which they appear.
  characters           Lists the characters in the book
  display              Displays the processed structure
  serve                Starts booktools as a webservice

Flags:
  -r, --chapterRegex string   Regular expression which if matched on a line will trigger a chapter.
  -h, --help                  help for process

Global Flags:
      --config string   config file (default is $HOME/.temp.yaml)
```

### Example

```
> ./booktools process serve mybook.txt
To access booktools, open a webbrowser and
navigate to http://localhost:8080/


```

![Example of character matches](http://drive.google.com/uc?id=1ZamoAaztjehdJF5uv2YPkgTD4_Eqyd4I)
