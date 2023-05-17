package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"encoding/json"
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

func createStringWriter(csvPath string) func(string, bool) {
	jsonDir := filepath.Dir(csvPath)
	jsonName := fmt.Sprintf("%s.json", strings.TrimSuffix(filepath.Base(csvPath), ".csv"))
	finalLocation := fmt.Sprintf("%s/%s", jsonDir, jsonName)

	f, err := os.Create(finalLocation)
	check(err)

	return func(data string, close bool) {
		_, err := f.WriteString(data)
		check(err)

		if close {
			f.Close()
		}
	}
}

func getJSONFunc(pretty bool) (func(map[string]string) string, string) {
	var jsonFunc func(map[string]string) string
	var breakLine string
	if pretty {
		breakLine = "\n"
		jsonFunc = func(record map[string]string) string {
			jsonData, _ := json.MarshalIndent(record, "   ", "   ")
			return "   " + string(jsonData)
		}
	} else {
		breakLine = ""
		jsonFunc = func(record map[string]string) string {
			jsonData, _ := json.Marshal(record)
			return string(jsonData)
		}
	}

	return jsonFunc, breakLine
}

func writeJSONFile(csvPath string, writerChannel <-chan map[string]string, done chan<- bool, pretty bool) {
	writeString := createStringWriter(csvPath)
	jsonFunc, breakLine := getJSONFunc(pretty)

	fmt.Println("Writing JSON file...")

	writeString("["+breakLine, false)
	first := true
	for {
		record, more := <-writerChannel
		if more {
			if !first {
				writeString(","+breakLine, false)
			} else {
				first = false
			}

			jsonData := jsonFunc(record)
			writeString(jsonData, false)
		} else {
			writeString(breakLine+"]", true)
			fmt.Println("Completed!")
			done <- true
			break
		}
	}
}

func main(){
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <csvFile>\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	fileData, err := getFileData()

	if err != nil {
		exit(err)
	}

	if _, err := checkIfValidFile(fileData.filePath); err != nil {
		exit(err)
	}

	writerChannel := make(chan map[string]string)
	done := make(chan bool)

	go processCsv(fileData, writerChannel)
	go writeJSONFile(fileData.filePath, writerChannel, done, fileData.pretty)

	<-done
}