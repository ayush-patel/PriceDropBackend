package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//Item as the basic data structure
type Item struct {
	Name          string  `json:"title"`
	Brand         string  `json:"description"`
	URL           string  `json:"url"`
	ID            int     `json:"id"`
	OriginalPrice float64 `json:"originalprice"`
	CurrentPrice  float64 `json:"currentprice"`
}

//Use a map as a temporary database
var itemStore = make(map[string]Item)

//Variable to generate key for the collection
var id int

//PostURLHandler /api/items
func PostURLHandler(w http.ResponseWriter, r *http.Request) {

	var item Item
	// Decode the incoming Note json
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		panic(err)
	}

	//Use a python scraper to get price at the posted URL
	cmd := exec.Command("python", "-c", "import priceFetch; print priceFetch.priceFetch('"+item.URL+"')")
	fetchedPriceBytes, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	fetchedPrice := strings.TrimSpace(string(fetchedPriceBytes[:]))
	if fetchedPrice == "Error" {
		panic(err)
	}

	price, err := strconv.ParseFloat(fetchedPrice, 64)
	if err != nil {
		panic(err)
	}

	id++
	item.ID = id //not a good implementation, but works for demo
	idString := strconv.Itoa(item.ID)

	item.CurrentPrice = price
	item.OriginalPrice = price

	item.Name = "Name " + idString   //TODO: update it
	item.Brand = "Brand " + idString //TODO: update it

	itemStore[idString] = item

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
	var items []Item
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
	// r.HandleFunc("/api/notes/{id}", PutURLHandler).Methods("PUT")
	// r.HandleFunc("/api/notes/{id}", DeleteItemHandler).Methods("DELETE")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	server := &http.Server{
		Addr:    port,
		Handler: r,
	}
	log.Println("Listening on" + port)
	server.ListenAndServe()
}
