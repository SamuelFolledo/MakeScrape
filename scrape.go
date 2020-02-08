package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	//SELECTORS LIST
	//Spotlight Deal = #refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-summary-card.card.ebayui-dne-item-featured-card--topDeals.ebayui-dne-featured-with-carousel > div > div > div.dne-itemtile-detail
	//Trending Deals = #refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-carousel.filmstrip-centered.ebayui-dne-carousel.ebayui-dne-trending-widget.filmstrip-1 > div
	//Trending Deals 2 (li:nth-child(9)) = #refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-carousel.filmstrip-centered.ebayui-dne-carousel.ebayui-dne-trending-widget.filmstrip-1 > div > ul > li:nth-child(9) > div > div.dne-itemtile-detail
	//Featured Deals: #refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-item-featured-card.ebayui-dne-item-featured-card > div
	//Featured Deals 2 (div:nth-child(1)): #refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-item-featured-card.ebayui-dne-item-featured-card > div > div:nth-child(1) > div > div.dne-itemtile-detail

	// On every a element which has href attribute call callback
	c.OnHTML("#refit-spf-container > div.sections-container > div.ebayui-dne-featured-card.ebayui-dne-featured-with-padding > div.ebayui-dne-item-featured-card.ebayui-dne-item-featured-card > div", func(e *colly.HTMLElement) {
		// link := e.Attr("href")
		// Print link
		print("\n++++++++E TEXT = ", e.Text, "\n\n")
		fmt.Printf("\n++++++++Link found: -> %s\n\n", e.Attr("href"))
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
}
