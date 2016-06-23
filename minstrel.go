package main

import (
	"encoding/json"
	"errors"
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
	searchEngines map[string]search.EngineInterface
}

// GetEngine retreive an engine for its id
func (minstrel *Minstrel) GetEngine(id string) (search.EngineInterface, error) {
	engine, success := minstrel.searchEngines[id]
	if !success {
		return nil, errors.New("Can't find requested engine")
	}
	return engine, nil
}

// GetSearch handle the Get method on the search url
func (minstrel *Minstrel) GetSearch(writer http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	q := vars["q"]
	if len(q) == 0 {
		http.Error(writer, "Invalid query", http.StatusBadRequest)
		return
	}
	engineID := vars["service"]
	engine, err := minstrel.GetEngine(engineID)
	if err != nil {
		http.Error(writer, "Couldn't find requested service", http.StatusBadRequest)
		return
	}
	minstrel.Search(writer, engine, q)
}

// GetDetail handle the Get method on the detail url
func (minstrel *Minstrel) GetDetail(writer http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)
	ids := vars["ids"]
	if len(ids) == 0 {
		http.Error(writer, "Invalid query", http.StatusBadRequest)
		return
	}
	engineID := vars["service"]
	engine, err := minstrel.GetEngine(engineID)
	if err != nil {
		http.Error(writer, "Couldn't find requested service", http.StatusBadRequest)
		return
	}
	minstrel.Detail(writer, engine, strings.Split(ids, ","))
}

// Search writes a json array of TrackID to the writer
func (minstrel *Minstrel) Search(writer http.ResponseWriter, engine search.EngineInterface, q string) {
	ids, err := search.Search(engine, q, "FR")
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
func (minstrel *Minstrel) Detail(writer http.ResponseWriter, engine search.EngineInterface, ids []string) {
	tracks, err := search.Detail(engine, ids)
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
