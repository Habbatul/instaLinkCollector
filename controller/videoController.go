package controller

import (
	"encoding/json"
	"instaLinkCollector/adapter"
	"net/http"
)

type VideoController struct {
	VideoApiAdapter adapter.VideoApiAdapter
}

func NewVideoController(apiAdapter *adapter.VideoApiAdapter) *VideoController {
	return &VideoController{}
}

func (v *VideoController) ShowListResolution(w http.ResponseWriter, r *http.Request) {
	response, err := v.VideoApiAdapter.FetchInstagramWithCookies(r.URL.Query().Get("video_url"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
