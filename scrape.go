package main

import (
	"fmt"
	"strconv" //convert string to float
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"

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
	// testDB()
	scrapeEbay()
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

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
		colly.Async(true),                    // Turn on asynchronous requests
		colly.Debugger(&debug.LogDebugger{}), // Attach a debugger to the collector
	)

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	// defer db.Close()
	// Migrate the schema
	db.AutoMigrate(&Deal{})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*", // when visiting links which domains' matches "*httpbin.*" glob
		Parallelism: 2,            // Limit the number of threads started by colly to two
		//Delay:      5 * time.Second,
	})

	//SELECTORS LIST
	//Spotlight Deal = #refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-summary-card.card.ebayui-dne-item-featured-card--topDeals.ebayui-dne-featured-with-carousel > div > div > div.dne-itemtile-detail
	//Trending Deals = #refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-carousel.filmstrip-centered.ebayui-dne-carousel.ebayui-dne-trending-widget.filmstrip-1 > div
	//Trending Deals 2 (li:nth-child(9)) = #refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-carousel.filmstrip-centered.ebayui-dne-carousel.ebayui-dne-trending-widget.filmstrip-1 > div > ul > li:nth-child(9) > div > div.dne-itemtile-detail
	//Featured Deals: #refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-item-featured-card.ebayui-dne-item-featured-card > div
	//Featured Deals 2 (div:nth-child(1)): #refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-item-featured-card.ebayui-dne-item-featured-card > div > div:nth-child(1) > div > div.dne-itemtile-detail

	// On every a element which has href attribute call callback
	c.OnHTML("#refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-item-featured-card.ebayui-dne-item-featured-card > div > div > div > div.dne-itemtile-detail", func(e *colly.HTMLElement) {
		// link := e.Attr("href")
		// Print link
		print("\n==========================================================")
		// print("\n++++++++E TEXT = ", e.Text, "\n\n")
		// print("\n++++++++E DOM = ", e.DOM.Text(), "\n\n") //e.DOM.Text() returns the same as e.Text
		var dealsScraped = e.DOM.Text()
		deal := createDealFromScrapedText(dealsScraped)
		db.Create(&deal)
		// fmt.Printf("\n++++++++Link found: -> %s\n\n", e.Attr("href"))
		// fmt.Println("PARAGRAPHS", e.DOM.Find("p").Text())
		// var scrapeResult = e.DOM.Text()
		print("----------------------------------------------------------\n")
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.ebay.com/deals")

	// for i := 0; i < 1; i++ { // Start scraping in five threads on https://httpbin.org/delay/2
	// 	c.Visit(fmt.Sprintf("%s?n=%d", "https://www.ebay.com/deals", i))
	// }

	c.Wait() // Wait until threads are finished

	deals := Deal{}
	db.Find(&deals)

	fmt.Println("\n\nDEALS ARE ===== ", deals.Name)
	// for _, deal := range deals {

	// }

	defer db.Close()
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
