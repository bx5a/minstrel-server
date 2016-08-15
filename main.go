package main

import (
	"log"
	"net/http"

	"github.com/bx5a/minstrel-server/search"
	"github.com/gorilla/mux"
)

func main() {
	engine := search.MakeYoutubeEngine()
	minstrel := &Minstrel{searchEngine: engine}
	r := mux.NewRouter()

	r.Methods("GET").Path("/v1/TrackIDs").HandlerFunc(minstrel.GetTrackIDs)
	r.Methods("GET").Path("/v1/Tracks").HandlerFunc(minstrel.GetTracks)

	log.Fatal(http.ListenAndServe(":8080", r))
}
