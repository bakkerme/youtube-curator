package youtubeapi

import (
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/utils"
	"strings"
)

// APIRequester provides an interface for access to the Youtube API
type APIRequester interface {
	GetVideoMetadata(ids *[]string, cf *config.Config) (*VideoMetadataResponse, error)
	GetVideosForChannel(ytc collection.YTChannel, cf *config.Config) (*VideoMetadataResponse, error)
}

// API allows access to the Youtube API
type API struct{}

var apiVideos = "videos"
var apiSearch = "search"
var apiPlaylistItems = "playlistItems"

// GetVideoMetadata gets information on the video whos IDs are provided from the Youtube API
func (ytAPI *API) GetVideoMetadata(ids *[]string, cf *config.Config) (*VideoMetadataResponse, error) {
	return getVideoMetadata(ids, cf, &utils.HTTPClient{})
}

func getVideoMetadata(ids *[]string, cf *config.Config, httpClient utils.YTCHTTPClient) (*VideoMetadataResponse, error) {
	if len(*ids) > 50 {
		return nil, errors.New("YT API cannot get more than 50 IDs at a time")
	}

	values := map[string]string{
		"part":  "snippet",
		"id":    strings.Join(*ids, ","),
		"order": "date",
	}

	body, err := makeAPIRequest(apiVideos, &values, getAccessKey(cf), httpClient)
	if err != nil {
		return nil, fmt.Errorf("Could not get Video Metadata from Youtube API. Error %s", err)
	}

	videoResponse, err := convertAPIResponse(string(body), apiVideos)

	if err != nil {
		return nil, fmt.Errorf("Could not parse response from Youtube API for video ID %s. Responsed with %s. Error %s", ids, body, err)
	}

	return videoResponse, nil
}

// GetVideosForChannel returns videos for a provided YT Channel from the YouTube API
func (ytAPI *API) GetVideosForChannel(ytc collection.YTChannel, cf *config.Config) (*VideoMetadataResponse, error) {
	return getVideosForChannel(ytc.ChannelType(), ytc, cf, &utils.HTTPClient{})
}

func getVideosForChannel(channelType string, ytc collection.YTChannel, cf *config.Config, httpClient utils.YTCHTTPClient) (*VideoMetadataResponse, error) {
	values := map[string]string{
		"part":  "snippet",
		"order": "date",
	}

	api := ""
	if channelType == collection.ChannelTypeChannel {
		values["channelId"] = ytc.ID()
		values["type"] = "video"
		api = apiSearch
	} else if channelType == collection.ChannelTypePlaylist {
		values["playlistId"] = ytc.ID()
		api = apiPlaylistItems
	} else {
		return nil, fmt.Errorf("Invalid Channel Type provided. Got %s", channelType)
	}

	body, err := makeAPIRequest(api, &values, getAccessKey(cf), httpClient)
	if err != nil {
		return nil, fmt.Errorf("Could not get Video Metadata from Youtube API for channel %s. Error %s", ytc.ID(), err)
	}

	videoResponse, err := convertAPIResponse(string(body), api)

	if err != nil {
		return nil, fmt.Errorf("Could not parse response from Youtube API for chanel %s. Responsed with %s. Error %s", ytc.ID(), body, err)
	}

	return videoResponse, nil
}

var baseURL string = "https://youtube.googleapis.com/youtube/v3/"

func makeAPIRequest(api string, keyVals *map[string]string, accessKey string, httpClient utils.YTCHTTPClient) ([]byte, error) {
	if accessKey == "" {
		return nil, errors.New("A valid access key is required to make requests to the Youtube API")
	}

	queryParams := []string{
		"key=" + accessKey,
	}

	if keyVals != nil {
		for key, val := range *keyVals {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, val))
		}
	}

	baseURLWithAPI := fmt.Sprintf("%s%s?", baseURL, api)
	url := baseURLWithAPI + strings.Join(queryParams, "&")

	resp, body, err := httpClient.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Returned invalid response for URL %s. Error: %s", url, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Call for URL %s did not return 200. Returned %d. Body was %s", url, resp.StatusCode, body)
	}

	return body, nil
}

func getAccessKey(cf *config.Config) string {
	return cf.YoutubeAPIKey
}
