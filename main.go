package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// Operations
const (
	MULTIPLY = 0
	SUM      = 1
)

// Run with
//		go run .
// Send request with:
//		curl -F 'file=@/path/matrix.csv' "localhost:8080/echo"

func main() {

	// 1. Echo 
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {

		records := validate(w, r)

		if records != nil {
			echoArray(w, records)
		}

	})

	// 2. Invert 
	http.HandleFunc("/invert", func(w http.ResponseWriter, r *http.Request) {

		records := validate(w, r)

		if records != nil {
			invert(w, records)
		}

	})

	// 3. Flatten
	http.HandleFunc("/flatten", func(w http.ResponseWriter, r *http.Request) {

		records := validate(w, r)

		if records != nil {
			flatten(w, records)
		}

	})

	// 4. Sum
	http.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request) {

		records := validate(w, r)

		if records != nil {
			operationMath(w, records, SUM)
		}

	})

	// 5. Multiply
	http.HandleFunc("/multiply", func(w http.ResponseWriter, r *http.Request) {

		records := validate(w, r)

		if records != nil {
			operationMath(w, records, MULTIPLY)
		}

	})

	// OPEN THE SERVER TO RUN
	http.ListenAndServe(":8090", nil)

}

// Make sure that all data is correctly
func validate(w http.ResponseWriter, r *http.Request) [][]string {

	//Treat the panic error
	defer func() {
		if recover() != nil {
			writeMessage(w, "An error occurred while processing the data.")
		}
	}()

	file, multipartFileHeader, err := r.FormFile("file")

	// try open file
	if err != nil {
		writeMessageError(w, err, "Error when trying to open the file. Allow it to be included correctly!")
		return nil
	}

	defer file.Close()

	// Verify file extesion
	if !validExtension(multipartFileHeader.Filename) {
		writeMessage(w, "Invalid file to process. Please choose a .txt or .csv file")
		return nil
	}

	// Read file content
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		writeMessageError(w, err, "Invalid data entry.The data format is not as expected.\n\n "+
			"Example of valid data (without white spaces): \n"+
			"1,2,3 \n"+
			"4,5,6 \n"+
			"7,8,9")
		return nil
	}

	if !validData(records, w) {
		writeMessage(w, "Invalid data type entry. The matrix is not square. \n\n"+
			"Example of Square matrix (without white spaces): \n"+
			"1,2,3 \n"+
			"4,5,6 \n"+
			"7,8,9")

		return nil
	}

	return records // return the data
}

func echoArray(w http.ResponseWriter, data [][]string) {
	var response string

	for _, sla := range data {
		response = fmt.Sprintf("%s%s\n", response, strings.Join(sla, ","))
	}
	fmt.Fprint(w, response)
}

// invert values from array
func invert(w http.ResponseWriter, data [][]string) {

	finalString := ""

	arrayInverted := make([][]string, len(data))

	for i := 0; i < len(data); i++ {

		for j := 0; j < len(data[i]); j++ {

			arrayInverted[j] = append(arrayInverted[j], strings.TrimSpace(data[i][j]))

		}

	}

	for _, sla := range arrayInverted {
		finalString = fmt.Sprintf("%s%s\n", finalString, strings.Join(sla, ","))
	}

	writeMessage(w, finalString)
}

// put all matrix in 1 line
func flatten(w http.ResponseWriter, data [][]string) {

	finalString := ""

	for i := 0; i < len(data); i++ {

		for j := 0; j < len(data[i]); j++ {

			if (i == len(data)-1) && (j == len(data[i])-1) { //The last record
				finalString = finalString + data[i][j]
			} else {
				finalString = finalString + data[i][j] + ","
			}
		}
	}

	writeMessage(w, finalString)
}

// function resposible for MULTIPLY and SUM values
func operationMath(w http.ResponseWriter, data [][]string, operationType int) {

	finalValue := 0.0

	len := len(data)

	for i := 0; i < len; i++ {
		for _, value := range data[i] {
			f, _ := strconv.ParseFloat(value, 64)

			switch operationType {
			case MULTIPLY:

				if finalValue == 0.0 { // prevent multiplication by zero
					finalValue = 1.0
				}

				finalValue = finalValue * f // multiply values

			default:
				finalValue = finalValue + f // sum values
			}

		}
	}

	writeMessage(w, fmt.Sprintf("%v", finalValue))
}

func writeMessage(w http.ResponseWriter, message string) {
	w.Write([]byte(fmt.Sprintf(message)))
}
func writeMessageError(w http.ResponseWriter, err error, message string) {
	w.Write([]byte(fmt.Sprintf(message+"\n\nError message: %s", err.Error())))
}

// check if the matrix is square && integers numbers
func validData(data [][]string, w http.ResponseWriter) bool {

	len := len(data)
	var rows = 0
	var columnsDefault = 0
	var columns = 0

	for i := 0; i < len; i++ {
		rows++

		if i == 0 {
			columnsDefault = 0

			for _, value := range data[i] {

				columnsDefault++

				// check if have only integers numbers
				if _, error := strconv.Atoi(strings.TrimSpace(value)); error != nil {
					writeMessage(w, "For this operation, you must use numbers Integers only.\n\n")
					return false
				}
			}

		} else {

			columns = 0

			for _, value := range data[i] {

				columns++

				// check if have only integers numbers
				if _, error := strconv.Atoi(strings.TrimSpace(value)); error != nil {
					writeMessage(w, "For this operation, you must use numbers Integers only.\n\n")
					return false
				}
			}

			if columns != columnsDefault {
				return false
			}
		}

	}

	if rows == columnsDefault { // check if is square
		return true
	}

	return false
}

// check file extension
func validExtension(fileName string) bool {

	switch ext := filepath.Ext(fileName); strings.TrimSpace(ext) {
	case ".txt":
		return true
	case ".csv":
		return true
	case ".xlsx":
		return true
	default:
		return false
	}

}
