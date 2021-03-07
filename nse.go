package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/PuerkitoBio/goquery"
)

func main() {

	nfIndices := []string{"NIFTY%20BANK","NIFTY%20CONSUMER%20DURABLES","NIFTY%20AUTO",
	"NIFTY%20FIN%20SERVICE","NIFTY%20FINSRV25%2050","NIFTY%20FMCG",
	"NIFTY%20Healthcare%20Index","NIFTY%20IT","NIFTY%20MEDIA","NIFTY%20METAL",
	"NIFTY%20OIL%20%26%20GAS","NIFTY%20PHARMA","NIFTY%20PVT%20BANK",
	"NIFTY%20PSU%20BANK","NIFTY%20REALTY",
	}

	frDate := "02-12-2020"
	toDate := "22-02-2021"

	for _, v := range nfIndices {
		url := "https://www1.nseindia.com/products/dynaContent/equities/indices/historicalindices.jsp?indexType="+v+"&fromDate="+frDate+"&toDate="+toDate
		call(url, "GET", v)
	}
}

func call(url string, method string, name string) {
	var headings, row []string
	var rows [][]string

	client := &http.Client{Timeout: time.Second * 10,}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36 Edg/88.0.705.74")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Language", "en-US")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("StatusCode error: %v", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			s.Find("th").Each(func(i int, s *goquery.Selection) {
				headings = append(headings, s.Text())
			})
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				row = append(row, s.Text())
			})
			rows = append(rows, row)
			row = nil
		})
	})

	// Filter empty rows.
	var rows_filtered [][]string
	for i := 0 ; i < len(rows)-1 ; i++ {
		if rows[i] != nil {
			rows_filtered = append(rows_filtered, rows[i])
		}
	}

	// Heading
	head := headings[0]
	
	// Period
	period := headings[1]

	// Columns
	var col_map = map[string]string{}

	for i, c := range headings[2:] {
		col_map[string(65+i)+strconv.Itoa(3)] = c
	}

	// Rows
	var rows_map = map[string]string{}

	for i := 0; i < len(rows_filtered); i++ {
		r := strings.Join(rows_filtered[i], ",")
		e := strings.Split(r, ",")
		for cn, v := range e {
				rows_map[string(65+cn)+strconv.Itoa(4+i)] = strings.Trim(v, " ")
		}
	}

	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", head)
	f.SetCellValue("Sheet1", "A2", period)

	for k, v := range col_map {
		f.SetCellValue("Sheet1", k, v)
	}

	for k, v := range rows_map {
		f.SetCellValue("Sheet1", k, v)
	}

    // Save spreadsheet by the given path.
    if err := f.SaveAs(name + ".xlsx"); err != nil {
        fmt.Println(err)
    }
	
}
