package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/bx5a/minstrel-server/search"
	"github.com/bx5a/minstrel-server/track"
)

/*
Minstrel is the main app
*/
type Minstrel struct {
	searchEngine search.EngineInterface
}

// GetTrackIDs search for tracks
func (minstrel *Minstrel) GetTrackIDs(writer http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	// get page token. If none, an empty string is set
	pageToken := r.URL.Query().Get("pageToken")
	if len(q) != 0 {
		minstrel.Search(writer, q, pageToken)
		return
	}
	http.Error(writer, "Invalid request", http.StatusBadRequest)
}

// GetTracks search for track details for a given list of TrackID
func (minstrel *Minstrel) GetTracks(writer http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("ids")
	if len(ids) != 0 {
		decoder := json.NewDecoder(strings.NewReader(ids))
		tracks := []track.ID{}
		for {
			var id track.ID
			err := decoder.Decode(&id)
			if err != nil {
				break
			}
			tracks = append(tracks, id)
		}
		minstrel.Detail(writer, tracks)
		return
	}
	http.Error(writer, "Invalid request", http.StatusBadRequest)
}

// Search writes a json array of TrackID to the writer
func (minstrel *Minstrel) Search(writer http.ResponseWriter, q string, pageToken string) {
	idList, err := search.Search(minstrel.searchEngine, q, "US", pageToken)
	if err != nil {
		log.Fatal(err)
	}
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	err2 := json.NewEncoder(writer).Encode(idList)
	if err2 != nil {
		log.Fatal(err)
	}
}

// Detail writes a json array of Track to the writer using the give ids
func (minstrel *Minstrel) Detail(writer http.ResponseWriter, ids []track.ID) {
	tracks, err := search.Detail(minstrel.searchEngine, ids)
	if err != nil {
		log.Fatal(err)
	}
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
	err2 := json.NewEncoder(writer).Encode(tracks)
	if err2 != nil {
		log.Fatal(err)
	}
}
