package search

import (
	"testing"

	"github.com/bx5a/minstrel-server/track"
)

func TestYoutubeEngine_Search(t *testing.T) {
	searchEngine := MakeYoutubeEngine()
	idList, err := searchEngine.Search("adele", "US", "")
	if err != nil {
		t.Fatal(err)
	}
	ids := idList.IDs
	if len(ids) == 0 {
		t.Errorf("Obtained ID array is empty")
	}
	if idList.NextPageToken == "" {
		t.Errorf("Invalid next page token")
	}

	if testing.Short() {
		t.Skip("Skipping search result test in shot mode.")
		return
	}

	if ids[0].Source != youtubeEngineSourceName {
		t.Errorf("ids[0].Source == %s, want %s", ids[0].Source, youtubeEngineSourceName)
	}
	expectedID := "fk4BbF7B29w"
	if ids[0].ID != expectedID {
		t.Errorf("ids[0].ID == %s, want %s", ids[0].ID, expectedID)
	}
}

func TestYoutubeEngine_Detail(t *testing.T) {
	searchEngine := MakeYoutubeEngine()
	adeleFirstTrack := track.ID{ID: "fk4BbF7B29w", Source: youtubeEngineSourceName}
	adeleSecondTrack := track.ID{ID: "YQHsXMglC9A", Source: youtubeEngineSourceName}
	ids := []track.ID{adeleFirstTrack, adeleSecondTrack}
	tracks, err := searchEngine.Detail(ids)
	if err != nil {
		t.Fatal(err)
	}

	if testing.Short() {
		t.Skip("Skipping detail result test in shot mode.")
		return
	}

	if len(tracks) != len(ids) {
		t.Errorf("len(tracks) != len(ids)")
	}

	if tracks[0].ID.ID != ids[0].ID {
		t.Errorf("tracks[0].ID.ID != ids[0].ID")
	}

	expectedTitle := "Adele - Send My Love (To Your New Lover)"
	if tracks[0].Title != expectedTitle {
		t.Errorf("tracks[0].Title == %s, want %s", tracks[0].Title, expectedTitle)
	}

	if tracks[0].Thumbnail.Default == "" {
		t.Errorf("Thumnail default url empty")
	}

	expectedDuration := "226000"
	if tracks[0].Duration != expectedDuration {
		t.Errorf("Unexpected duration: %s, want %s", tracks[0].Duration, expectedDuration)
	}
}

func TestYoutubeEngine_Category(t *testing.T) {
	searchEngine := MakeYoutubeEngine()
	ID, err := searchEngine.queryMusicCategoryID()
	if err != nil {
		t.Fatal(err)
	}
	if ID == "" {
		t.Errorf("Obtained ID is empty")
	}
}
