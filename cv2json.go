package main

import(
	"errors"
	"flag"
	"os"
)

type inputFile struct{
	filePath string
	seperator string
	pretty bool
}

func getFileData() (inputFile, error){
	if len(os.Args) < 2 {
		return inputFile{}, errors.New("A filepath argument is required")
	}

	seperator := flag.String("seperator","comma","Column seperator")
	pretty := flag.Bool("pretty", false, "Generate pretty JSON")

	flag.Parse()

	fileLocation := os.Args(0)

	if !( *seperator == "comma" || *seperator == "semicolon" ){
		return inputFile{}, errors.New("Only comma or semicolon seperators allowed")
	}

	return inputFile{fileLocation, *seperator, *pretty}, nil
}

func main(){

}