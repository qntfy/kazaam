package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/qntfy/kazaam"
)

var (
	// command-line arguments
	specFilename = flag.String("spec", "", "Kazaam Specification (required)")
	inFilename   = flag.String("in", "", "Input file (optional)")
	outFilename  = flag.String("out", "", "Output file (optional)")
	verbose      = flag.Bool("verbose", false, "Turn on verbose logging")
)

func loadKazaamTransform(specFilename string) *kazaam.Kazaam {
	if specFilename == "" {
		log.Fatalln("Must specify a Kazaam specification file.")
	}
	specFile, specFileError := ioutil.ReadFile(specFilename)
	if specFileError != nil {
		log.Fatal("Unable to read Kazaam specification file: ", specFileError)
	}
	k, specError := kazaam.NewKazaam(string(specFile))
	if specError != nil {
		log.Fatal("Unable to load Kazaam specification file: ", specError)
	}
	return k
}

func main() {
	flag.Parse()

	k := loadKazaamTransform(*specFilename)

	var inData []byte
	var in, out string
	var readError error
	if *inFilename == "" {
		// read from stdin
		reader := bufio.NewReader(os.Stdin)
		inData, readError = ioutil.ReadAll(reader)
	} else {
		// read from file
		inData, readError = ioutil.ReadFile(*inFilename)
	}
	in = string(inData)
	if readError != nil {
		log.Fatal("Unable to read input", readError)
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
