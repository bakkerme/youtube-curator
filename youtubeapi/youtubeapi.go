package youtubeapi

import (
	"config"
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

// func getVideoInfo(id string, cf *config.Config) (string, error) {
// url := getVideoInfoURL(id, getAccessKey(cf))

// // return "", fmt.Errorf("Returned invalid response for address %s. Response was %d", url, resp.StatusCode)
// }
