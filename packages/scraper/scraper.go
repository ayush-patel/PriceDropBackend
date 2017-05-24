package scraper

import (
	"errors"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//scraper for nordstrom.com
func nordstronScrape(url string) (string, string, float64, string, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return "", "", 0, "", errors.New("Can't load the webpage")
	}

	productDetails := doc.Find(".product-information")
	name := productDetails.Find("h1").Text()
	if name == "" {
		return "", "", 0, "", errors.New("Can't retrieve the name")
	}

	price, _ := strconv.ParseFloat(productDetails.Find(".current-price").Text()[1:], 64)
	if price == 0 {
		return "", "", 0, "", errors.New("Can't retrieve price")
	}

	imageURL, _ := doc.Find(".thumbnail-wrapper .image-thumbnail a img").Attr("src")
	if imageURL == "" {
		return "", "", 0, "", errors.New("Can't retrieve image URL")
	}

	return "Nordstrom", name, price, imageURL, nil
}

//Scrape returns the scraped data (brand, name, price, imageURL, error)
func Scrape(url string) (string, string, float64, string, error) {
	if strings.Contains(url, "nordstrom") {
		return nordstronScrape(url)
	}
	return "", "", 0, "", errors.New("Webscraper is not available for the given site")
}
