package youtubeapi

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/utils"
	"strings"
)

func getAccessKey(cf *config.Config) string {
	return cf.YoutubeAPIKey
}

func getVideoInfoURL(id string, accessKey string) string {
	values := []string{
		"part=snippet",
		"id=" + id,
		"key=" + accessKey,
	}

	baseURL := "https://youtube.googleapis.com/youtube/v3/videos?"
	return baseURL + strings.Join(values, "&")
}

func getVideoInfo(id string, cf *config.Config, httpClient utils.YTCHTTPClient) (*VideoListResponse, error) {
	url := getVideoInfoURL(id, getAccessKey(cf))

	resp, body, err := httpClient.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Returned invalid response for video Id %s, URL was %s. Error: %s", id, url, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Call for video ID %s with URL %s did not return 200. Returned %d. Body was %s", id, url, resp.StatusCode, body)
	}

	videoResponse, err := convertVideoAPIResponse(string(body))

	if err != nil {
		return nil, fmt.Errorf("Could not parse response from Youtube API for video ID %s, address %s. Responsed with %s", id, url, body)
	}

	return videoResponse, nil
}

// GetVideoInfo gets the video information for the provided ID from the Youtube API
func GetVideoInfo(id string, cf *config.Config) (*VideoListResponse, error) {
	return getVideoInfo(id, cf, &utils.HTTPClient{})
}
