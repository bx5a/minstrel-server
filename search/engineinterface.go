package search

import "github.com/bx5a/minstrel-server/track"

/*
EngineInterface defines function required in a search engine
*/
type EngineInterface interface {
	Search(q string, countryCode string, pageToken string) (track.IDList, error)
	Detail(ids []track.ID) ([]track.Track, error)
}

// Search gives a unified way of searching accord Engine implementations
func Search(engine EngineInterface, q string, countryCode string, pageToken string) (track.IDList, error) {
	return engine.Search(q, countryCode, pageToken)
}

// Detail gives a unified way of getting details accord Engine implementations
func Detail(engine EngineInterface, ids []track.ID) ([]track.Track, error) {
	return engine.Detail(ids)
}
