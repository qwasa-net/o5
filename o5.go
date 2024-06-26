package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// flagsSet is all command-line options in one structure
type flagsSet struct {
	inFileName  string
	outFileName string
	tokenTrim   bool
	macStart    string
	macEnd      string
	workDir     string
	defines     map[string]string
	definesFile string
}

// main is an entry point
func main() {
	flags := readFlags()
	inFile, outFile := getFiles(flags)
	defer inFile.Close()
	defer outFile.Close()
	parse(inFile, outFile, flags)
}

// parse reads from inFile and writes to outFile
func parse(inFile io.Reader, outFile io.Writer, flags flagsSet) {

	var err error

	// macro expression
	qStart, qEnd := regexp.QuoteMeta(flags.macStart), regexp.QuoteMeta(flags.macEnd)
	macXP := regexp.MustCompile(qStart + "(.+?)" + qEnd)

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

// boombastia expands tokens with file content or variable
func boombastia(token string, flags flagsSet) string {

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

// getFiles reads flags and returns Files to read from and write to
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

// check panics
func check(err error) {
	if err != nil {
		panic(err)
	}
}

type flagsArray []string

func (i *flagsArray) String() string         { return "" }
func (i *flagsArray) Set(value string) error { *i = append(*i, value); return nil }

// readFlags reads command line parameters
func readFlags() flagsSet {

	var flags = getDefaultFlags()

	flag.BoolVar(&flags.tokenTrim, "trim", flags.tokenTrim, "trim spaces in expanded macro")
	flag.StringVar(&flags.macStart, "start", flags.macStart, "macro openner (prefix)")
	flag.StringVar(&flags.macEnd, "end", flags.macEnd, "macro closer (suffix)")
	flag.StringVar(&flags.inFileName, "i", flags.inFileName, "input file ('-' is stdin)")
	flag.StringVar(&flags.outFileName, "o", flags.outFileName, "output file ('-' is stdout)")
	flag.StringVar(&flags.workDir, "w", flags.workDir, "working directory (for file includes)")
	var defines flagsArray
	flag.Var(&defines, "d", "define macro variable (-d NAME=VALUE)")
	flag.StringVar(&flags.definesFile, "dd", flags.definesFile, "read macro variables from file")

	usage := func() {
		exename := filepath.Base(os.Args[0])
		fmt.Fprint(os.Stderr, exename, " -- super simple micro macro processor for text files\n\n")
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

	if flags.definesFile != "" {
		data, err := os.ReadFile(flags.definesFile)
		check(err)
		for _, line := range strings.Split(string(data), "\n") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				flags.defines[parts[0]] = parts[1]
			} else {
				flags.defines[parts[0]] = ""
			}
		}
	}

	return flags

}

// getDefaultFlags creates flagsSet with default presets
func getDefaultFlags() flagsSet {

	return flagsSet{
		tokenTrim:   true,
		macStart:    "<!--#",
		macEnd:      "-->",
		inFileName:  "-",
		outFileName: "-",
		workDir:     ".",
		defines:     make(map[string]string),
	}

}
