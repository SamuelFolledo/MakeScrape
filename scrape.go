package main

import (
	"fmt"
	"net/http"
	"strconv" //convert string to float
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/labstack/echo"

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

	// db, err := gorm.Open("sqlite3", "test.db")
	// if err != nil {
	// 	panic("failed to connect database")
	// }
	// defer db.Close()

	// var deal Deal
	// db.First(&deal, 1)
	// println("DEAL 1 = ", deal.Name)
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
	// Instantiate default collector
	c := colly.NewCollector(
		colly.Async(false),                   // Turn on asynchronous requests, this stuff could return before or after.
		colly.Debugger(&debug.LogDebugger{}), // Attach a debugger to the collector
	)

	// db, err := gorm.Open("sqlite3", "test.db")
	// if err != nil {
	// 	panic("failed to connect database")
	// }
	// defer db.Close()
	// Migrate the schema
	// db.AutoMigrate(&Deal{})

	// c.Limit(&colly.LimitRule{ //limit must be use to async = true. replace httpbin to ebay.com
	// 	DomainGlob:  "*httpbin.*", // when visiting links which domains' matches "*httpbin.*" glob
	// 	Parallelism: 2,            // Limit the number of threads started by colly to two
	// 	//Delay:      5 * time.Second,
	// })

	var deals []string

	// On every a element which has href attribute call callback
	c.OnHTML("#refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-item-featured-card.ebayui-dne-item-featured-card > div > div > div > div.dne-itemtile-detail", func(e *colly.HTMLElement) {
		var dealsScraped = e.DOM.Text()
		deals = append(deals, dealsScraped)
		// deal := createDealFromScrapedText(dealsScraped)
		// db.Create(&deal)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	e := echo.New() //create a server, inir xollt
	e.GET("/scrape", func(ec echo.Context) (err error) {
		c.Visit("https://www.ebay.com/deals")
		c.Wait()

		var firstDeal string
		for _, item := range deals {
			firstDeal += "- " + item + "\n"
		}
		// return ec.JSON()
		return ec.String(http.StatusOK, firstDeal)
	})

	e.Logger.Fatal(e.Start(":1323"))

	c.Visit("https://www.ebay.com/deals")

	c.Wait() // Wait until threads are finished
}

func createDealFromScrapedText(dealsScraped string) Deal {
	dealArray := strings.Split(dealsScraped, "$") //separate string by dollar sign
	var name = dealArray[0]
	currentPrice, _ := strconv.ParseFloat(dealArray[1], 8)
	var previousPrice float64
	if len(dealArray) > 2 { //if array's length is > 2, then we have previousPrice
		prev, _ := strconv.ParseFloat(dealArray[2], 8)
		previousPrice = prev
	}
	// for index, word := range s {
	print("\nDEAL IS ", name, " CURRENT PRICE = ", currentPrice, " PREVIOUSLY = ", previousPrice)
	return Deal{Name: name, CurrentPrice: currentPrice, PreviousPrice: previousPrice}
}
