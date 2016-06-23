package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/bx5a/minstrel-server/search"
	"github.com/gorilla/mux"
)

/*
Minstrel is the main app
*/
type Minstrel struct {
	searchEngine search.EngineInterface
}

// GetSearch handle the Get method on the search url
func (minstrel *Minstrel) GetSearch(writer http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	q := vars["q"]
	if len(q) == 0 {
		http.Error(writer, "'q' component missing", http.StatusBadRequest)
		return
	}
	minstrel.Search(writer, q)
}

// GetDetail handle the Get method on the detail url
func (minstrel *Minstrel) GetDetail(writer http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids := vars["ids"]
	if len(ids) == 0 {
		http.Error(writer, "'ids' component missing", http.StatusBadRequest)
		return
	}
	minstrel.Detail(writer, strings.Split(ids, ","))
}

// Search writes a json array of TrackID to the writer
func (minstrel *Minstrel) Search(writer http.ResponseWriter, q string) {
	ids, err := search.Search(minstrel.searchEngine, q, "FR")
	if err != nil {
		log.Fatal(err)
	}

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	err2 := json.NewEncoder(writer).Encode(ids)
	if err2 != nil {
		log.Fatal(err)
	}
}

// Detail writes a json array of Track to the writer using the give ids
func (minstrel *Minstrel) Detail(writer http.ResponseWriter, ids []string) {
	tracks, err := search.Detail(minstrel.searchEngine, ids)
	if err != nil {
		log.Fatal(err)
	}

	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	err2 := json.NewEncoder(writer).Encode(tracks)
	if err2 != nil {
		log.Fatal(err)
	}
}
