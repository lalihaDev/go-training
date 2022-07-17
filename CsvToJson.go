package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	csvPath := flag.String("csvPath", ".", "Path of the CSV files directory")
	jsonPath := flag.String("jsonPath", ".", "Path of the JSon files directory")
	flag.Parse()

	csvDirectory, err := os.Open(*csvPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	files, err := csvDirectory.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {

		absolutePath, err := filepath.Abs(file.Name())
		if err != nil {
			fmt.Println(err)
			return
		}
		//TODO fix error :  The file is not found || wrong root
		fileBytes, filePath := ReadCSV(absolutePath, jsonPath)
		SaveFile(fileBytes, filePath)
		fmt.Println(file.Name())
	}

	fmt.Println("Csv files are converted to JSon with success")
}

// ReadCSV to read the content of CSV File
func ReadCSV(csvPath string, jsonPath *string) ([]byte, string) {
	csvFile, err := os.OpenFile(csvPath, os.O_RDWR, 0644)

	if err != nil {
		log.Fatal("The file is not found || wrong root")
	}

	//TODO check the file length : it should have less then 1 mega

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	content, _ := reader.ReadAll()

	if len(content) < 1 {
		log.Fatal("Something wrong, the file maybe empty or length of the lines are not the same")
	}

	headersArr := retrieveCsvHeaders(content)

	//Remove the header row
	content = content[1:]

	return buildJSon(jsonPath, content, headersArr)
}

//buildJSon to build content of JSon file from csv data: headers and content
//the return value is the JSon file
//TODO to refactor : maybe exists a library doing this??!!
func buildJSon(path *string, content [][]string, headersArr []string) ([]byte, string) {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for i, data := range content {
		buffer.WriteString("{")
		for j, y := range data {
			buffer.WriteString(`"` + headersArr[j] + `":`)
			_, fErr := strconv.ParseFloat(y, 32)
			_, bErr := strconv.ParseBool(y)
			if fErr == nil {
				buffer.WriteString(y)
			} else if bErr == nil {
				buffer.WriteString(strings.ToLower(y))
			} else {
				buffer.WriteString((`"` + y + `"`))
			}
			//end of property
			if j < len(data)-1 {
				buffer.WriteString(",")
			}

		}
		//end of object of the array
		buffer.WriteString("}")
		if i < len(content)-1 {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString(`]`)
	rawMessage := json.RawMessage(buffer.String())
	jsonFile, _ := json.MarshalIndent(rawMessage, "", "  ")
	jsonFileName := filepath.Base(*path)
	jsonFileName = jsonFileName[0:len(jsonFileName)-len(filepath.Ext(jsonFileName))] + ".json"
	jsonDirectory := filepath.Dir(*path)
	return jsonFile, filepath.Join(jsonDirectory, jsonFileName)
}

//retrieveCsvHeaders extract headers from csv file
//the return value is an array of csv headers
func retrieveCsvHeaders(content [][]string) []string {
	headersArr := make([]string, 0)
	for _, headE := range content[0] {
		headersArr = append(headersArr, headE)
	}
	return headersArr
}

// SaveFile Will Save the file
func SaveFile(myFile []byte, path string) {
	if err := ioutil.WriteFile(path, myFile, os.FileMode(0644)); err != nil {
		panic(err)
	}
}
