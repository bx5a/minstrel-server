package search

import (
	"net/http"
	"strings"

	"github.com/bx5a/minstrel-server/track"

	"github.com/google/google-api-go-client/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// TODO: paste your api key here
const developerKey = "AIzaSyAOdGv1-HH3nNPsn-EBzFSpSO5IcTkKtMI"
const sourceName = "YouTube"

/*
YoutubeEngine is the EngineInterface implementation for youtube
*/
type YoutubeEngine struct {
}

// Search is an inherited function from EngineInterface
func (engine YoutubeEngine) Search(q string, countryCode string) ([]track.ID, error) {
	client := &http.Client{Transport: &transport.APIKey{Key: developerKey}}
	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	call := service.Search.List("id").
		Q(q).
		RegionCode(countryCode).
		Type("video").
		MaxResults(50)

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	// Iterate through each item and add it to the list
	ids := []track.ID{}
	for _, item := range response.Items {
		ids = append(ids, track.ID{ID: item.Id.VideoId, Source: sourceName})
	}
	return ids, nil
}

// Detail returns the detail for the list of track from a list of youtube ids
func (engine YoutubeEngine) Detail(ids []string) ([]track.Track, error) {
	client := &http.Client{Transport: &transport.APIKey{Key: developerKey}}
	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	call := service.Videos.List("snippet,contentDetails").
		Id(strings.Join(ids[:], ","))

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	tracks := []track.Track{}
	for index, item := range response.Items {
		trackID := track.ID{ID: ids[index], Source: sourceName}
		track := track.Track{ID: trackID, Title: item.Snippet.Title, Duration: item.ContentDetails.Duration}
		tracks = append(tracks, track)
	}
	return tracks, nil
}
