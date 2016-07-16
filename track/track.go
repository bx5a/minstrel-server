package track

/*
Track is a detailed representation of something that can be played
*/
type Track struct {
	ID        ID        `json:"id"`
	Title     string    `json:"title"`
	Duration  string    `json:"duration"`
	Thumbnail Thumbnail `json:"thumbnail"`
}
