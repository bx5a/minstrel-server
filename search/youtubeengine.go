package search

import (
	"errors"
	"net/http"
	"strings"

	"github.com/bx5a/minstrel-server/track"

	"github.com/google/google-api-go-client/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

/*
YoutubeEngine is the EngineInterface implementation for youtube
*/
type YoutubeEngine struct {
}

// Search is an inherited function from EngineInterface
func (engine YoutubeEngine) Search(q string, countryCode string) ([]track.ID, error) {
	client := &http.Client{Transport: &transport.APIKey{Key: youtubeEngineDeveloperKey}}
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
		ids = append(ids, track.ID{ID: item.Id.VideoId, Source: youtubeEngineSourceName})
	}
	return ids, nil
}

// Detail returns the detail for the list of track from a list of youtube ids
func (engine YoutubeEngine) Detail(ids []track.ID) ([]track.Track, error) {
	// if any of the id requested has an invalid source, return an error
	stringIds := []string{}
	for _, element := range ids {
		if element.Source != youtubeEngineSourceName {
			return nil, errors.New("Invalid source required")
		}
		stringIds = append(stringIds, element.ID)
	}

	client := &http.Client{Transport: &transport.APIKey{Key: youtubeEngineDeveloperKey}}
	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	call := service.Videos.List("snippet,contentDetails").
		Id(strings.Join(stringIds, ","))

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	tracks := []track.Track{}
	for index, item := range response.Items {
		track := track.Track{
			ID:       ids[index],
			Title:    item.Snippet.Title,
			Duration: item.ContentDetails.Duration,
			Thumbnail: track.Thumbnail{
				Default: item.Snippet.Thumbnails.Default.Url,
				High:    item.Snippet.Thumbnails.High.Url,
			},
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}
