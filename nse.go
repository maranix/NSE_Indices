package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://www1.nseindia.com/products/dynaContent/equities/indices/historicalindices.jsp?indexType=NIFTY%20BANK&fromDate=02-12-2020&toDate=22-02-2021"
	call(url, "GET")
}

func call(url string, method string) {
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

	// fmt.Println("####### headings = ", len(headings), headings)

	// f := excelize.NewFile()
	// f.SetCellValue("Sheet1", "A1", headings[0])
    // // Save spreadsheet by the given path.
    // if err := f.SaveAs("Book1.xlsx"); err != nil {
    //     fmt.Println(err)
    // }

	fmt.Println("####### rows = ", len(rows), rows)
	

	// bytes, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(string(bytes))
}