package main

import (
	"log"
	"net/http"

	"github.com/bx5a/minstrel-server/search"
	"github.com/gorilla/mux"
)

func main() {
	minstrel := &Minstrel{searchEngine: search.YoutubeEngine{}}
	r := mux.NewRouter()

	r.Methods("GET").Path("/v1/search/{q}").HandlerFunc(minstrel.GetSearch)
	r.Methods("GET").Path("/v1/detail/youtube/{ids}").HandlerFunc(minstrel.GetDetail)

	log.Fatal(http.ListenAndServe(":8080", r))
}
