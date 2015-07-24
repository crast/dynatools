/*
update-reserved is a tiny utility which generates the reserved words list for safeupdate by parsing amazon's site.
*/
package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const ReservedListUrl = "http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/ReservedWords.html"
const OutFile = "reserved.go"

func main() {
	doc, err := goquery.NewDocument(ReservedListUrl)
	fatalErr(err)
	reserved := doc.Find("pre.programlisting").First().Text()

	fp, err := os.Create(OutFile)
	if os.IsExist(err) {
		fp, err = os.OpenFile(OutFile, os.O_RDWR, 0666)
	}
	fatalErr(err)
	buf := bytes.NewBufferString(baseFile)

	for _, s := range strings.Split(strings.TrimSpace(reserved), "\n") {
		log.Print(s)
		buf.WriteString(strconv.Quote(s) + ": true,\n")
	}
	buf.WriteString("\n}")
	b, err := format.Source(buf.Bytes())
	fatalErr(err)
	_, err = fp.Write(b)
	fatalErr(err)

}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var baseFile = `package safeupdate
// DO NOT EDIT THIS FILE, it's been auto-generated
//
//go:generate go run update-reserved/main.go

var reservedWords = map[string]bool{
`