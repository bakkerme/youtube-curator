package youtubeapi

import (
	"encoding/json"
	"fmt"
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

// ResourceID specifies the type of Resource in a result
type ResourceID struct {
	Kind string `json:"kind,omitempty"`
	ID   string `json:"videoId,omitempty"`
}

// VideoSnippet Contains general information about the video
type VideoSnippet struct {
	PublishedAt          string             `json:"publishedAt,omitempty"`
	PublishTime          string             `json:"publishTime,omitempty"`
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
	PlaylistID           string             `json:"playlistId,omitempty"`
	Position             int                `json:"position,omitempty"`
	ResourceID           ResourceID         `json:"resourceId,omitempty"`
}

// Video Represents a single YouTube video
type Video struct {
	Kind    string       `json:"kind,omitempty"`
	Etag    string       `json:"etag,omitempty"`
	ID      string       `json:"id,omitempty"`
	Snippet VideoSnippet `json:"snippet,omitempty"`
}

// SearchVideo Represents a single YouTube video
type SearchVideo struct {
	Kind    string       `json:"kind,omitempty"`
	Etag    string       `json:"etag,omitempty"`
	ID      ResourceID   `json:"id,omitempty"`
	Snippet VideoSnippet `json:"snippet,omitempty"`
}

// PlaylistItemMetadataResponse The top level return from the video list API
type PlaylistItemMetadataResponse struct {
	Kind          string   `json:"kind,omitempty"`
	Etag          string   `json:"etag,omitempty"`
	NextPageToken string   `json:"nextPageToken,omitempty"`
	Items         []Video  `json:"items,omitempty"`
	PageInfo      PageInfo `json:"pageInfo,omitempty"`
}

// SearchMetadataResponse The top level return from the video list API
type SearchMetadataResponse struct {
	Kind          string        `json:"kind,omitempty"`
	Etag          string        `json:"etag,omitempty"`
	NextPageToken string        `json:"nextPageToken,omitempty"`
	RegionCode    string        `json:"regionCode,omitempty"`
	Items         []SearchVideo `json:"items,omitempty"`
	PageInfo      PageInfo      `json:"pageInfo,omitempty"`
}

// VideoMetadataResponse The top level return from the video list API
type VideoMetadataResponse struct {
	Kind          string   `json:"kind,omitempty"`
	Etag          string   `json:"etag,omitempty"`
	NextPageToken string   `json:"nextPageToken,omitempty"`
	RegionCode    string   `json:"regionCode,omitempty"`
	Items         []Video  `json:"items,omitempty"`
	PageInfo      PageInfo `json:"pageInfo,omitempty"`
}

func convertAPIResponse(file string, api string) (*VideoMetadataResponse, error) {
	switch api {
	case apiSearch:
		return convertSearchJSONToVideo(file)
	case apiPlaylistItems:
		return convertPlaylistItemJSONToVideo(file)
	case apiVideos:
		return convertVideoAPIResponse(file)
	default:
		return nil, fmt.Errorf("%s is not a valid API type", api)
	}
}

func convertVideoAPIResponse(file string) (*VideoMetadataResponse, error) {
	var resp VideoMetadataResponse
	if err := json.Unmarshal([]byte(file), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func convertSearchJSONToVideo(file string) (*VideoMetadataResponse, error) {
	var resp SearchMetadataResponse
	if err := json.Unmarshal([]byte(file), &resp); err != nil {
		return nil, err
	}

	videoItems := []Video{}
	for _, item := range resp.Items {
		convertedVideo := Video{
			Kind:    item.Kind,
			Etag:    item.Etag,
			ID:      item.ID.ID,
			Snippet: item.Snippet,
		}

		videoItems = append(videoItems, convertedVideo)
	}

	vmr := VideoMetadataResponse{
		Kind:          resp.Kind,
		Etag:          resp.Etag,
		NextPageToken: resp.NextPageToken,
		RegionCode:    resp.RegionCode,
		PageInfo:      resp.PageInfo,
		Items:         videoItems,
	}

	return &vmr, nil
}

func convertPlaylistItemJSONToVideo(file string) (*VideoMetadataResponse, error) {
	var resp PlaylistItemMetadataResponse
	if err := json.Unmarshal([]byte(file), &resp); err != nil {
		return nil, err
	}

	videoItems := []Video{}
	for _, item := range resp.Items {
		convertedVideo := Video{
			Kind:    item.Kind,
			Etag:    item.Etag,
			ID:      item.Snippet.ResourceID.ID,
			Snippet: item.Snippet,
		}

		videoItems = append(videoItems, convertedVideo)
	}

	vmr := VideoMetadataResponse{
		Kind:          resp.Kind,
		Etag:          resp.Etag,
		NextPageToken: resp.NextPageToken,
		PageInfo:      resp.PageInfo,
		Items:         videoItems,
	}

	return &vmr, nil
}
