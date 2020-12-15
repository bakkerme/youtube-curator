package youtubeapi

import (
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/utils"
	"strings"
)

// GetVideoInfo gets information on the video whos IDs are provided from the Youtube API
func GetVideoInfo(ids *[]string, cf *config.Config) (*VideoListResponse, error) {
	return getVideoInfo(ids, cf, &utils.HTTPClient{})
}

func getAccessKey(cf *config.Config) string {
	return cf.YoutubeAPIKey
}

func getVideoInfoURL(ids *[]string, accessKey string) string {
	values := []string{
		"part=snippet",
		"id=" + strings.Join(*ids, ","),
		"key=" + accessKey,
	}

	baseURL := "https://youtube.googleapis.com/youtube/v3/videos?"
	return baseURL + strings.Join(values, "&")
}

func getVideoInfo(ids *[]string, cf *config.Config, httpClient utils.YTCHTTPClient) (*VideoListResponse, error) {
	if len(*ids) > 50 {
		return nil, errors.New("YT API cannot get more than 50 IDs at a time")
	}

	url := getVideoInfoURL(ids, getAccessKey(cf))

	resp, body, err := httpClient.Get(url)

	if err != nil {
		return nil, fmt.Errorf("Returned invalid response for video Id %s, URL was %s. Error: %s", ids, url, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Call for video ID %s with URL %s did not return 200. Returned %d. Body was %s", ids, url, resp.StatusCode, body)
	}

	videoResponse, err := convertVideoAPIResponse(string(body))

	if err != nil {
		return nil, fmt.Errorf("Could not parse response from Youtube API for video ID %s, address %s. Responsed with %s", ids, url, body)
	}

	return videoResponse, nil
}
