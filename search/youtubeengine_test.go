package search

import (
	"testing"

	"github.com/bx5a/minstrel-server/track"
)

func TestYoutubeEngine_Search(t *testing.T) {
	searchEngine := YoutubeEngine{}
	ids, err := searchEngine.Search("adele", "US")
	if err != nil {
		t.Fatal(err)
	}

	if testing.Short() {
		t.Skip("Skipping search result test in shot mode.")
		return
	}

	if ids[0].Source != sourceName {
		t.Errorf("ids[0].Source == %s, want %s", ids[0].Source, sourceName)
	}
	expectedID := "fk4BbF7B29w"
	if ids[0].ID != expectedID {
		t.Errorf("ids[0].ID == %s, want %s", ids[0].ID, expectedID)
	}
}

func TestYoutubeEngine_Detail(t *testing.T) {
	searchEngine := YoutubeEngine{}
	adeleFirstTrack := track.ID{ID: "fk4BbF7B29w", Source: sourceName}
	adeleSecondTrack := track.ID{ID: "YQHsXMglC9A", Source: sourceName}
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
}