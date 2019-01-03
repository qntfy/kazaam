// A simple command-line interface (CLI) for executing kazaam transforms on data from files or stdin.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mbordner/kazaam"
)

var (
	// command-line arguments
	specFilename = flag.String("spec", "", "Kazaam Specification (required)")
	inFilename   = flag.String("in", "", "Input file (optional)")
	outFilename  = flag.String("out", "", "Output file (optional)")
	verbose      = flag.Bool("verbose", false, "Turn on verbose logging")
)

func loadKazaamTransform(specFilename string) (*kazaam.Kazaam, error) {
	if specFilename == "" {
		return nil, errors.New("Must specify a Kazaam specification file")
	}
	specFile, specFileError := ioutil.ReadFile(specFilename)
	if specFileError != nil {
		return nil, errors.New("Unable to read Kazaam specification file: " + specFileError.Error())
	}
	k, specError := kazaam.NewKazaam(string(specFile))
	if specError != nil {
		return nil, errors.New("Unable to load Kazaam specification file: " + specError.Error())
	}
	return k, nil
}

func getInput(inputFilename string, file *os.File) (string, error) {
	var inData []byte
	var readError error
	if inputFilename == "" {
		// read from stdin
		reader := bufio.NewReader(file)
		inData, readError = ioutil.ReadAll(reader)
	} else {
		// read from file
		inData, readError = ioutil.ReadFile(inputFilename)
	}
	if readError != nil {
		return "", readError
	}
	return string(inData), nil
}

func main() {
	flag.Parse()

	k, err := loadKazaamTransform(*specFilename)
	if err != nil {
		log.Fatal("Trouble loading specification", err)
	}

	in, err := getInput(*inFilename, os.Stdin)
	if err != nil {
		log.Fatal("Unable to open specified input")
	}

	out, transformError := k.TransformJSONStringToString(in)
	if transformError != nil {
		log.Fatal("Unable to transform message", transformError)
	}

	if *outFilename == "" {
		// write to stdout
		fmt.Print(out)
	} else {
		// write to file
		writeError := ioutil.WriteFile(*outFilename, []byte(out), 0644)
		if writeError != nil {
			log.Fatal("Unable to write transformed output", writeError)
		}
	}
}
