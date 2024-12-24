package dto

type VideoVersion struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
	Type   int    `json:"type"`
}
