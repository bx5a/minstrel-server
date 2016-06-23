package main

import (
	"log"
	"net/http"

	"github.com/bx5a/minstrel-server/search"
	"github.com/gorilla/mux"
)

func main() {
	var engines map[string]search.EngineInterface
	engines = make(map[string]search.EngineInterface)

	engines["youtube"] = search.YoutubeEngine{}
	minstrel := &Minstrel{searchEngines: engines}
	r := mux.NewRouter()

	r.Methods("GET").Path("/v1/search/{service}/{q}").HandlerFunc(minstrel.GetSearch)
	r.Methods("GET").Path("/v1/detail/{service}/{ids}").HandlerFunc(minstrel.GetDetail)

	log.Fatal(http.ListenAndServe(":8080", r))
}
