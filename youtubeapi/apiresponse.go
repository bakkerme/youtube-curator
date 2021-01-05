package youtubeapi

import (
	"encoding/json"
)

// PageInfo Information on the pagination of the API request
type PageInfo struct {
	TotalResults   int `json:"totalResults,omitempty"`
	ResultsPerPage int `json:"resultsPerPage,omitempty"`
}

// LocalizationDetail Localized video data
type LocalizationDetail struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// Thumbnail An individual thumbnail of a video
type Thumbnail struct {
	URL    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// ThumbnailDetails Data on the various thumbnail sizes available
type ThumbnailDetails struct {
	Default  Thumbnail `json:"default,omitempty"`
	Medium   Thumbnail `json:"medium,omitempty"`
	High     Thumbnail `json:"high,omitempty"`
	Standard Thumbnail `json:"standard,omitempty"`
	Maxres   Thumbnail `json:"maxres,omitempty"`
}

// VideoSnippet Contains general information about the video
type VideoSnippet struct {
	PublishedAt          string             `json:"publishedAt,omitempty"`
	ChannelID            string             `json:"channelId,omitempty"`
	Title                string             `json:"title,omitempty"`
	Description          string             `json:"description,omitempty"`
	Thumbnails           ThumbnailDetails   `json:"thumbnails,omitempty"`
	ChannelTitle         string             `json:"channelTitle,omitempty"`
	Tags                 []string           `json:"tags,omitempty"`
	CategoryID           string             `json:"categoryId,omitempty"`
	LiveBroadcastContent string             `json:"liveBroadcastContent,omitempty"`
	DefaultLanguage      string             `json:"defaultLanguage,omitempty"`
	Localized            LocalizationDetail `json:"localized,omitempty"`
	DefaultAudioLanguage string             `json:"defaultAudioLanguage,omitempty"`
}

// Video Represents a single YouTube video
type Video struct {
	Kind    string       `json:"kind,omitempty"`
	Etag    string       `json:"etag,omitempty"`
	ID      string       `json:"id,omitempty"`
	Snippet VideoSnippet `json:"snippet,omitempty"`
}

// VideoMetadataResponse The top level return from the video list API
type VideoMetadataResponse struct {
	Kind          string   `json:"kind,omitempty"`
	Etag          string   `json:"etag,omitempty"`
	NextPageToken string   `json:"nextPageToken,omitempty"`
	RegionCode    string   `json:"regionCode,omitempty"`
	PageInfo      PageInfo `json:"pageInfo,omitempty"`
	Items         []Video  `json:"items,omitempty"`
}

func convertVideoAPIResponse(file string) (*VideoMetadataResponse, error) {
	var resp VideoMetadataResponse
	if err := json.Unmarshal([]byte(file), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
