package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type empData struct {
	Product string
	Price   int
	Rating  int
}

type fileInterface interface {
	FileFormat()
}

type fileName struct {
	Name string
}

func (fn fileName) FileFormat() {
	// check file format and choice correct func
	fileType := strings.Split(fn.Name, ".")
	print(fileType[1])
	if fileType[1] == "csv" {
		csvReader(fn.Name)
	} else if fileType[1] == "json" {
		jsonReader(fn.Name)
	}
}

func main() {
	fmt.Println("Search...")

	var fi fileInterface

	fn := fileName{Name: os.Args[1]}
	fi = fn
	fi.FileFormat()

	fmt.Println("done")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func csvReader(filename string) {
	csvFile, err := os.Open("./data/" + filename)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	var maxRatingRow []*empData
	var maxPriceRow []*empData

	for id, line := range csvLines {
		if id == 0 {
			continue // pass first row ("Product,Price,Rating")
		}
		priceInt, err := strconv.Atoi(line[1])
		checkErr(err)
		ratingInt, err := strconv.Atoi(line[2])
		checkErr(err)

		if id == 1 {
			maxPriceRow = append(maxPriceRow, &empData{line[0], priceInt, ratingInt})
			maxRatingRow = append(maxRatingRow, &empData{line[0], priceInt, ratingInt})
		}
		price := *maxPriceRow[0]
		rait := *maxRatingRow[0]
		emp := &empData{
			Product: line[0],
			Price:   priceInt,
			Rating:  ratingInt,
		}
		if emp.Price > price.Price {
			maxPriceRow[0] = emp
		}
		if emp.Rating > rait.Rating {
			maxRatingRow[0] = emp
		}
	}
	rait := *maxRatingRow[0]
	price := *maxPriceRow[0]
	fmt.Println("max Price -", price, "<|>", "max Rating -", rait)
}

func jsonReader(filename string) {
	jsonFile, err := os.Open("./data/" + filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	checkErr(err)
	var result []empData

	jsonErr := json.Unmarshal(data, &result)
	checkErr(jsonErr)

	var maxRatingRow []*empData
	var maxPriceRow []*empData
	for k, v := range result {
		if k == 0 {
			maxPriceRow = append(maxPriceRow, &empData{v.Product, v.Price, v.Rating})
			maxRatingRow = append(maxRatingRow, &empData{v.Product, v.Price, v.Rating})
			continue
		}
		price := *maxPriceRow[0]
		rait := *maxRatingRow[0]

		emp := &empData{
			Product: v.Product,
			Price:   v.Price,
			Rating:  v.Rating,
		}
		if v.Price > price.Price {
			maxPriceRow[0] = emp
		}
		if v.Rating > rait.Rating {
			maxRatingRow[0] = emp
		}
	}
	rait := *maxRatingRow[0]
	price := *maxPriceRow[0]
	fmt.Println("max Price -", price, "<|>", "max Rating -", rait)
}
