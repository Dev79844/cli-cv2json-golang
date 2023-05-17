package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type inputFile struct{
	filePath string
	seperator string
	pretty bool
}

func getFileData() (inputFile, error) {
	if len(os.Args) < 2 {
		return inputFile{}, errors.New("a filepath argument is required")
	}

	separator := flag.String("separator", "comma", "Column separator")
	pretty := flag.Bool("pretty", false, "Generate pretty JSON")

	flag.Parse()

	fileLocation := flag.Arg(0)

	if !(*separator == "comma" || *separator == "semicolon") {
		return inputFile{}, errors.New("only comma or semicolon separators are allowed")
	}

	return inputFile{fileLocation, *separator, *pretty}, nil
}

func checkIfValidFile(filename string) (bool,error){

	if fileExtension := filepath.Ext(filename); fileExtension != ".csv"{
		return false, fmt.Errorf("file %s is not a csv",filename)
	}

	if _,err := os.Stat(filename); err != nil && os.IsNotExist(err){
		return false,fmt.Errorf("file %s does not exist",filename)
	}

	return true, nil
}


func main(){
	fileData,err := getFileData()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fileData)
}