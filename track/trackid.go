package track

/*
TrackID is the smallest object that can represent a track
*/
type ID struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}
