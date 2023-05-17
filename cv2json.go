package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
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

func exit(err error){
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
   	os.Exit(1)
}

func check(err error){
	if err != nil {
		check(err)
	}
}

func processLine(headers []string, dataList []string) (map[string]string, error) {
	if len(dataList) != len(headers) {
		return nil, errors.New("line doesn't match headers format. Skipping")
	}
	recordMap := make(map[string]string)
	for i, name := range headers {
		recordMap[name] = dataList[i]
	}
	return recordMap, nil
}

func processCsv(fileData inputFile, writerChannel chan<-map[string] string){
	file,err := os.Open(fileData.filePath)

	check(err)

	defer file.Close()

	var headers, line []string

	reader := csv.NewReader(file)

	if fileData.seperator == "semicolon"{
		reader.Comma = ';'
	}

	headers, err = reader.Read()

	check(err)

	for{
		line,err = reader.Read()

		if err == io.EOF{
			close(writerChannel)
			break
		}else if err != nil {
			exit(err)
		}

		record,err := processLine(headers, line)
		if err != nil {
			fmt.Printf("Line: %sError: %s\n", line, err)
			continue
		}

		writerChannel<-record

	}

}

func main(){
	fileData,err := getFileData()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fileData)
}