package youtubeapi

import (
	"config"
	"strings"
	"testing"
)

func TestGetAccessKey(t *testing.T) {
	t.Run("getAccessKey returns valid results", func(t *testing.T) {
		apiTestString := "123abc"

		key := getAccessKey(&config.Config{
			YoutubeAPIKey: apiTestString,
		})

		if key != apiTestString {
			t.Errorf("getAccessKey did not return correct results. Expected %s got %s", apiTestString, key)
		}
	})
}

func TestGetVideoInfoURL(t *testing.T) {
	t.Run("getVideoInfoURL returns a valid URL", func(t *testing.T) {
		id := "123abc"
		accessKey := "cba321"

		url := getVideoInfoURL(id, accessKey)

		hasHTTPS := strings.Contains(url, "https://")

		if !hasHTTPS {
			t.Errorf("video URL does not have http:// protocol: %s", url)
		}

		hasID := strings.Contains(url, id)
		if !hasID {
			t.Errorf("video URL does not have provided ID. Expected to find: %s, got: %s", id, url)
		}

		hasAccessKey := strings.Contains(url, accessKey)
		if !hasAccessKey {
			t.Errorf("video URL does not have provided access key. Expected to find: %s, got: %s ", accessKey, url)
		}
	})
}
