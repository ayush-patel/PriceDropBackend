package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/user/PriceDropBackend/packages/scraper"
)

//Item as the basic data structure
// type Item struct {
// 	Name          string  `json:"title"`
// 	Brand         string  `json:"description"`
// 	URL           string  `json:"url"`
// 	ImageURL      string  `json:"imageurl"`
// 	ID            int     `json:"id"`
// 	OriginalPrice float64 `json:"originalprice"`
// 	CurrentPrice  float64 `json:"currentprice"`
// }

//Use a map as a temporary database
var itemStore = make(map[string]scraper.Item)

//Variable to generate key for the collection
var id int = 1

var itemTemp scraper.Item

func dataGenerator(id int, name string, brand string, oprice float64, cprice float64, url string, iurl string) {
	itemTemp.ID = id
	itemTemp.Name = name
	itemTemp.Brand = brand
	itemTemp.OriginalPrice = oprice
	itemTemp.CurrentPrice = cprice
	itemTemp.URL = url
	itemTemp.ImageURL = iurl

	itemStore[strconv.Itoa(itemTemp.ID)] = itemTemp
}

//PostURLHandler /api/items
func PostURLHandler(w http.ResponseWriter, r *http.Request) {

	var item scraper.Item
	// Decode the incoming Note json
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		panic(err)
	}

	brand, name, price, imageURL, err := scraper.Scrape(item.URL)
	if err != nil {
		panic(err)
	}

	id++
	item.ID = id //not a good implementation, but works for demo
	item.CurrentPrice = price
	item.OriginalPrice = price
	item.Name = name
	item.Brand = brand
	item.ImageURL = imageURL

	itemStore[strconv.Itoa(item.ID)] = item

	j, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

//GetItemsHandler - /api/items
func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	//waitgroup to wait for the goroutines to get over when fetching prices
	var wg sync.WaitGroup

	var items []scraper.Item
	for _, v := range itemStore {
		wg.Add(1)
		go scraper.FetchPrice(&wg, &v)
	}

	wg.Wait() //wait for the goroutines to update prices before appending them
	for _, v := range itemStore {
		items = append(items, v)
	}

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(items)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//
// //HTTP Put - /api/notes/{id}
// func PutNoteHandler(w http.ResponseWriter, r *http.Request) {
// 	var err error
// 	vars := mux.Vars(r)
// 	k := vars["id"]
// 	var noteToUpd Note
// 	// Decode the incoming Note json
// 	err = json.NewDecoder(r.Body).Decode(&noteToUpd)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if note, ok := noteStore[k]; ok {
// 		noteToUpd.CreatedOn = note.CreatedOn
// 		//delete existing item and add the updated item
// 		delete(noteStore, k)
// 		noteStore[k] = noteToUpd
// 	} else {
// 		log.Printf("Could not find key of Note %s to delete", k)
// 	}
// 	w.WriteHeader(http.StatusNoContent)
// }
//
// //HTTP Delete - /api/notes/{id}
// func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	k := vars["id"]
// 	// Remove from Store
// 	if _, ok := noteStore[k]; ok {
// 		//delete existing item
// 		delete(noteStore, k)
// 	} else {
// 		log.Printf("Could not find key of Note %s to delete", k)
// 	}
// 	w.WriteHeader(http.StatusNoContent)
// }

//Entry point of the program
func main() {

	dataGenerator(0, "Impulse High Waist Midi Leggings", "Nordstrom", 65, 65, "http://shop.nordstrom.com/s/zella-impulse-high-waist-midi-leggings/4400095?origin=category-personalizedsort&fashioncolor=WHITE", "http://n.nordstrommedia.com/ImageGallery/store/product/Zoom/12/_100523132.jpg?crop=pad&pad_color=FFF&format=jpeg&w=60&h=90")
	dataGenerator(1, "Teagen Sneaker", "Nordstrom", 89.95, 89.95, "http://shop.nordstrom.com/s/halogen-teagen-sneaker-women/4528050?origin=category-personalizedsort&fashioncolor=BLACK%20LEATHER%20PERF", "http://n.nordstrommedia.com/ImageGallery/store/product/Zoom/4/_100715004.jpg?crop=pad&pad_color=FFF&format=jpeg&trim=color&trimcolor=FFF&w=60&h=90")
	dataGenerator(2, "Windsor Blazer", "Nordstrom", 495, 495, "http://shop.nordstrom.com/s/rag-bone-windsor-blazer/4214264?origin=category-personalizedsort&fashioncolor=BLACK", "http://n.nordstrommedia.com/ImageGallery/store/product/Zoom/2/_12023002.jpg?crop=pad&pad_color=FFF&format=jpeg&w=60&h=90")
	dataGenerator(3, "Gemini Link Tote", "Nordstrom", 195, 195, "http://shop.nordstrom.com/s/tory-burch-gemini-link-tote/4490562?origin=category-personalizedsort&fashioncolor", "http://n.nordstrommedia.com/ImageGallery/store/product/Zoom/0/_13480640.jpg?crop=pad&pad_color=FFF&format=jpeg&trim=color&trimcolor=FFF&w=60&h=90")
	dataGenerator(4, "Cover-Up Tunic", "Nordstrom", 62, 62, "http://shop.nordstrom.com/s/surf-gypsy-cover-up-tunic/4607134?origin=category-personalizedsort&fashioncolor=WHITE%2F%20CORAL%2F%20BLUE", "http://n.nordstrommedia.com/ImageGallery/store/product/Zoom/19/_100920639.jpg?crop=pad&pad_color=FFF&format=jpeg&w=60&h=90")
	dataGenerator(5, "'Leo' Envelope Clutch", "Nordstrom", 95, 95, "http://shop.nordstrom.com/s/rebecca-minkoff-leo-envelope-clutch/3853690?origin=category-personalizedsort&fashioncolor=OPTIC%20WHITE", "http://n.nordstrommedia.com/ImageGallery/store/product/Zoom/12/_100287932.jpg?crop=pad&pad_color=FFF&format=jpeg&trim=color&trimcolor=FFF&w=60&h=90")

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/items", GetItemsHandler).Methods("GET")
	r.HandleFunc("/api/items", PostURLHandler).Methods("POST")
	r.PathPrefix("/websites/nordstrom/").Handler(http.StripPrefix("/websites/nordstrom/", http.FileServer(http.Dir("./websites/nordstrom/"))))
	// r.HandleFunc("/api/notes/{id}", PutURLHandler).Methods("PUT")
	// r.HandleFunc("/api/notes/{id}", DeleteItemHandler).Methods("DELETE")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	log.Println("Listening on " + port)
	server.ListenAndServe()
}
