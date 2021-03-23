package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// command-line options in one structure
type flagsSet struct {
	inFileName  string
	outFileName string
	tokenTrim   bool
	macStart    string
	macEnd      string
	workDir     string
	defines     map[string]string
}

func main() {

	var err error

	flags := readFlags()
	inFile, outFile := getFiles(flags)

	defer inFile.Close()
	defer outFile.Close()

	// macro expression
	macXP := regexp.MustCompile(regexp.QuoteMeta(flags.macStart) + "(.+?)" + regexp.QuoteMeta(flags.macEnd))

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
				otoken := boombastia(token, flags)
				if flags.tokenTrim {
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

func boombastia(token string, flags flagsSet) string {
	// Expand TOKEN -- load a file, get variable

	// load file
	if strings.HasPrefix(token, "@") {
		data, err := ioutil.ReadFile(flags.workDir + token[1:])
		check(err)
		return string(data)
	}

	// get ENV variable
	if strings.HasPrefix(token, "$") {
		return os.Getenv(token[1:])
	}

	// get define
	return flags.defines[token]

}

func getFiles(flags flagsSet) (*os.File, *os.File) {

	var err error
	var inFile *os.File
	var outFile *os.File

	if flags.inFileName == "-" || flags.inFileName == "" {
		inFile = os.Stdin
	} else {
		inFile, err = os.Open(flags.inFileName)
		check(err)
	}

	if flags.outFileName == "-" || flags.outFileName == "" {
		outFile = os.Stdout
	} else {
		outFile, err = os.Create(flags.outFileName)
		check(err)
	}

	return inFile, outFile

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type flagsArray []string

func (i *flagsArray) String() string         { return "" }
func (i *flagsArray) Set(value string) error { *i = append(*i, value); return nil }

func readFlags() flagsSet {

	// read Command-Line Flags

	var flags = flagsSet{}
	flags.defines = make(map[string]string)

	flag.BoolVar(&flags.tokenTrim, "trim", true, "trim spaces in expanded macro")
	flag.StringVar(&flags.macStart, "start", "<!--#", "macro openner (prefix)")
	flag.StringVar(&flags.macEnd, "end", "-->", "macro closer (suffix)")
	flag.StringVar(&flags.inFileName, "i", "-", "input file ('-' is stdin)")
	flag.StringVar(&flags.outFileName, "o", "-", "output file ('-' is stdout)")
	flag.StringVar(&flags.workDir, "w", ".", "working directory (for file includes)")

	var defines flagsArray
	flag.Var(&defines, "d", "define macro variable (-d NAME=VALUE)")

	usage := func() {
		exename := filepath.Base(os.Args[0])
		fmt.Println(exename, `-- super simple micro macro processor for text files`)
		flag.PrintDefaults()
	}

	flag.Usage = usage

	flag.Parse()

	if len(defines) > 0 {
		for _, def := range defines {
			parts := strings.Split(def, "=")
			if len(parts) == 2 {
				flags.defines[parts[0]] = parts[1]
			}
		}
	}

	return flags

}
