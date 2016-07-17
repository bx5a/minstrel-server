package track

/*
IDList defines a list of track ID
*/
type IDList struct {
	IDs           []ID   `json: "tracks"`
	NextPageToken string `json: "nextPageToken"`
}
