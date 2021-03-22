package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

//..<!--# hallo -->*::*<!--# bye-bye -->..//
//..file=<!--# @include.txt -->=..user=<!--# USER -->=
func main() {

	var err error
	var inFile *os.File
	var outFile *os.File

	// 2do move to args
	inFileName := "-"
	outFileName := "-"
	tokenTrim := true
	macStart := "<!--#"
	macEnd := "-->"
	macXP := regexp.MustCompile(macStart + "(.+?)" + macEnd)

	if inFileName == "-" || inFileName == "" {
		inFile = os.Stdin
	} else {
		inFile, err = os.Open(inFileName)
		check(err)
		defer inFile.Close()
	}

	if outFileName == "-" || outFileName == "" {
		outFile = os.Stdout
	} else {
		outFile, err = os.Create(outFileName)
		check(err)
		defer inFile.Close()
	}

	scanner := bufio.NewScanner(inFile)
	writer := bufio.NewWriter(outFile)
	for scanner.Scan() {
		var itext, otext, token string
		itext = scanner.Text()
		matches := macXP.FindAllStringSubmatchIndex(itext, -1)
		if len(matches) > 0 {
			otext = ""
			i := 0
			for _, v := range matches {
				// match = [start, end, g1_start, g2_end]
				otext += itext[i:v[0]]
				token = strings.Trim(itext[v[2]:v[3]], " \n\t")
				otoken := replace(token)
				if tokenTrim {
					otoken = strings.Trim(otoken, " \n\t")
				}
				otext += otoken
				i = v[1]
			}
			otext += itext[i:]
		} else {
			otext = itext
		}
		writer.WriteString(otext + "\n")
	}
	writer.Flush()

	err = scanner.Err()
	check(err)

}

func replace(input string) string {

	// load file
	if strings.HasPrefix(input, "@") {
		data, err := ioutil.ReadFile(input[1:])
		check(err)
		return string(data)
	}

	// return ENV variable
	return os.Getenv(input)

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
