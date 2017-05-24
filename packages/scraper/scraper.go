package scraper

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

//scraper for nordstrom.com
func nordstromFetchData(url string) (string, string, float64, string, error) {
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

func nordstromFetchPrice(url string) (float64, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return 0, errors.New("Can't load the webpage")
	}

	price, _ := strconv.ParseFloat(doc.Find(".product-information .current-price").Text()[1:], 64)
	if price == 0 {
		return 0, errors.New("Can't retrieve price")
	}

	return price, nil
}

//Scrape returns the scraped data (brand, name, price, imageURL, error)
func Scrape(url string) (string, string, float64, string, error) {
	if strings.Contains(url, "nordstrom") {
		return nordstromFetchData(url)
	}
	return "", "", 0, "", errors.New("Webscraper is not available for the given site")
}

//Item as the basic data structure
type Item struct {
	Name          string  `json:"title"`
	Brand         string  `json:"description"`
	URL           string  `json:"url"`
	ImageURL      string  `json:"imageurl"`
	ID            int     `json:"id"`
	OriginalPrice float64 `json:"originalprice"`
	CurrentPrice  float64 `json:"currentprice"`
}

//FetchPrice if the price is changed the bool is returned as true
func FetchPrice(wg *sync.WaitGroup, item *Item) (bool, error) {
	defer wg.Done() //wait until after the process to mark it as done

	if strings.Contains(item.URL, "nordstrom") {
		fetchedPrice, err := nordstromFetchPrice(item.URL)
		if err != nil {
			return false, err
		}
		if fetchedPrice < item.CurrentPrice {
			item.CurrentPrice = fetchedPrice
			return true, nil
		}
		return false, nil
	}
	return false, errors.New("Webscraper is not available for the given site")
}
