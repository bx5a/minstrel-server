package search

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/bx5a/minstrel-server/track"

	"github.com/ChannelMeter/iso8601duration"
	"github.com/google/google-api-go-client/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

/*
YoutubeEngine is the EngineInterface implementation for youtube
*/
type YoutubeEngine struct {
	httpClient      *http.Client
	musicCategoryID string
}

// MakeYoutubeEngine builds a YouTube search engine. constructor pattern
func MakeYoutubeEngine() YoutubeEngine {
	client := &http.Client{Transport: &transport.APIKey{Key: youtubeEngineDeveloperKey}}
	return YoutubeEngine{httpClient: client, musicCategoryID: ""}
}

// musicCategoryId returns the YouTube id of the music category
func (engine YoutubeEngine) queryMusicCategoryID() (string, error) {
	service, err := youtube.New(engine.httpClient)
	if err != nil {
		return "", err
	}

	// TODO: Region code shouldn't always be US. It's not a problem in that case because
	// we query the music category but can be a problem in other calls
	call := service.VideoCategories.List("id,snippet").
		RegionCode("US").
		Fields("items(id,snippet/title)")
	response, err := call.Do()

	if err != nil {
		return "", err
	}

	for _, item := range response.Items {
		// TODO: make it a constant
		if item.Snippet.Title == "Music" {
			return item.Id, nil
		}
	}
	return "", errors.New("Music ID couldn't be found")
}

// getMusicCategoryID initialize and return the music category ID
func (engine YoutubeEngine) getMusicCategoryID() string {
	if engine.musicCategoryID == "" {
		ID, err := engine.queryMusicCategoryID()
		if err != nil {
			return ""
		}
		engine.musicCategoryID = ID
	}
	return engine.musicCategoryID
}

// Search is an inherited function from EngineInterface
func (engine YoutubeEngine) Search(q string, countryCode string, pageToken string) (track.IDList, error) {
	idList := track.IDList{IDs: nil, NextPageToken: ""}
	service, err := youtube.New(engine.httpClient)
	if err != nil {
		return idList, err
	}

	call := service.Search.List("id").
		Q(q).
		RegionCode(countryCode).
		VideoCategoryId(engine.getMusicCategoryID()).
		Type("video").
		PageToken(pageToken)

	response, err := call.Do()
	if err != nil {
		return idList, err
	}

	// Iterate through each item and add it to the list
	ids := []track.ID{}
	for _, item := range response.Items {
		ids = append(ids, track.ID{ID: item.Id.VideoId, Source: youtubeEngineSourceName})
	}

	idList.IDs = ids
	idList.NextPageToken = response.NextPageToken
	return idList, nil
}

// ISO8601DurationToMilliseconds converts a ISO 8601 duration to milliseconds
func (engine YoutubeEngine) ISO8601DurationToMilliseconds(isoDuration string) int {
	res, err := duration.FromString(isoDuration)
	if err != nil {
		return 0
	}
	weekPerYear := 52.1429
	dayPerWeek := 7.0
	hourPerDay := 24.0
	minPerHour := 60.0
	secPerMin := 60.0
	millisecPerSec := 1000.0

	weeks := float64(res.Years)*weekPerYear + float64(res.Weeks)
	days := weeks*dayPerWeek + float64(res.Days)
	hours := days*hourPerDay + float64(res.Hours)
	min := hours*minPerHour + float64(res.Minutes)
	sec := min*secPerMin + float64(res.Seconds)

	return int(sec * millisecPerSec)
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

	service, err := youtube.New(engine.httpClient)
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
		milliseconds := engine.ISO8601DurationToMilliseconds(item.ContentDetails.Duration)
		track := track.Track{
			ID:       ids[index],
			Title:    item.Snippet.Title,
			Duration: strconv.Itoa(milliseconds),
			Thumbnail: track.Thumbnail{
				Default: item.Snippet.Thumbnails.Default.Url,
				High:    item.Snippet.Thumbnails.High.Url,
			},
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}
