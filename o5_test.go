package main

import (
	"bytes"
	"strings"
	"testing"
)

// Test0 - does nothing
func Test0(t *testing.T) {}

// TestInput —
func TestInput(t *testing.T) {

	var input, output, want string

	input = "hello, $# there #$ $# 2+2 #$"
	want = "hello, here 4"

	inFile := bytes.NewBuffer([]byte(input))
	outFile := bytes.NewBuffer([]byte(""))

	var flags = getDefaultFlags()
	flags.defines["there"] = "here"
	flags.defines["2+2"] = "4"
	flags.macStart = "$#"
	flags.macEnd = "#$"

	parse(inFile, outFile, flags)

	output = outFile.String()
	output = strings.Trim(output, " \n\t")

	if output != want {
		t.Fatalf(`Test("%q") = %q, want %q`, input, output, want)
	}

}
