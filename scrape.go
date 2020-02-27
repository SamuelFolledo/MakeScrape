package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv" //convert string to float
	"strings"

	//colly - colleting/web scraper
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"

	//echo - making it live
	"github.com/labstack/echo"
	//gorm - database CRUD
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Deal struct {
	gorm.Model
	Name          string  `json:"name" gorm:"name"`
	CurrentPrice  float64 `json:"currentPrice" gorm:"currentPrice"`
	PreviousPrice float64 `json:"previousPrice" gorm:"previousPrice"`
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	scrapeEbay()
}

func testDB() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Deal{})

	// Create
	db.Create(
		&Deal{Name: "Deal 1", CurrentPrice: 1000, PreviousPrice: 9999},
	)

	// Read
	var deal Deal
	db.First(&deal, 1)                    // find product with id 1
	db.First(&deal, "name = ?", "Deal 1") // find product with code L1212

	fmt.Println("Deal = ", deal.Name, " ", deal.CurrentPrice)

	// Update - update product's price to 2000
	db.Model(&deal).Update("Price", 2000)

	fmt.Println("Deal = ", deal.Name, " ", deal.CurrentPrice)

	// Delete - delete product
	db.Delete(&deal)
}

func scrapeEbay() {

	c := colly.NewCollector( // Instantiate default collector
		colly.Async(false),                   // Turn on asynchronous requests, this stuff could return before or after.
		colly.Debugger(&debug.LogDebugger{}), // Attach a debugger to the collector
	)

	var deals [][]byte //a json is type []byte, so [][]byte creates a slice/array of []byte

	// On every a element which has href attribute call callback
	c.OnHTML("#refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-item-featured-card.ebayui-dne-item-featured-card > div > div > div > div.dne-itemtile-detail", func(e *colly.HTMLElement) {
		var dealsScraped = e.DOM.Text()
		deal := createDealFromScrapedText(dealsScraped)       //create a Deal from scraped text
		dealJson, err := json.MarshalIndent(deal, "", "    ") //makes it more pretty printed
		if isError(err) {
			return
		}
		deals = append(deals, dealJson) //append deal created to deal json
		// db.Create(&deal)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit("https://www.ebay.com/deals")

	c.Wait() // Wait until threads are finished

	//Having this here doesn't work, but if created from e.GET then it works
	jsonString := dealSliceToString(deals) //create a string from deals array
	writeToFile("output.json", jsonString) //write it to a json file

	//Start echo and display it on our browser
	e := echo.New() //create a server
	e.GET("/scrape", func(ec echo.Context) (err error) {
		return ec.String(http.StatusOK, jsonString) //display jsonString to echo
	})

	e.Logger.Fatal(e.Start(":1323"))
}

//creates a Deal instance from a string
func createDealFromScrapedText(dealsScraped string) Deal {
	dealArray := strings.Split(dealsScraped, "$") //separate string by dollar sign
	var name string
	var currentPrice float64 //ParseFloat converts the string s to a floating-point number with the precision specified by bitSize: 32 for float32, or 64 for float64.
	var previousPrice float64
	if len(dealArray) > 0 {
		name = dealArray[0]
	}
	if len(dealArray) > 1 {
		current, _ := strconv.ParseFloat(dealArray[1], 8)
		currentPrice = current
	}
	if len(dealArray) > 2 { //if array's length is > 2, then we have previousPrice
		prev, _ := strconv.ParseFloat(dealArray[2], 8)
		previousPrice = prev
	}
	print("\nDEAL IS ", name, " CURRENT PRICE = ", currentPrice, " PREVIOUSLY = ", previousPrice)
	return Deal{Name: name, CurrentPrice: currentPrice, PreviousPrice: previousPrice}
}

func isError(err error) bool { //error helper
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	return (err != nil)
}

//write to a file given a name and lines to write
func writeToFile(fileName, lines string) {
	bytesToWrite := []byte(lines)                         //data written
	err := ioutil.WriteFile(fileName, bytesToWrite, 0644) //filename, byte array (binary representation), and 0644 which represents permission number. (0-777) //will create a new text file if that text file does not exist yet
	if isError(err) {
		return
	}
}

//function that takes a deal slice/array and returns it in a pretty print json string format
func dealSliceToString(deals [][]byte) string {
	var jsonString string
	for _, dealJson := range deals { //loop through each JSON and add it to jsonString
		jsonString += (string(dealJson) + ",\n") //add a comma and a new line each dealJson
		// var res Deal
		// json.Unmarshal(bytes, &res) //returns it back to Deal struct
	}
	return jsonString
}

// Calculate returns x + 2.
func Calculate(x int) (result int) {
	result = x + 2
	return result
}
