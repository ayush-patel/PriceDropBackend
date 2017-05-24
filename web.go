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
var id int

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
