package config

import (
	"fmt"
	"testing"
)

type TestingConfigProvider struct{}

var ytTestKey string = "123abc"

func (cp TestingConfigProvider) getValue(key string) (string, error) {
	switch key {
	case "YOUTUBE_API_KEY":
		return ytTestKey, nil
	default:
		return "", fmt.Errorf("Config asked for a value not in the TestingConfigProvider. Wanted %s", key)
	}
}

func TestGetConfig(t *testing.T) {
	t.Run("Returns a valid config object for a provider", func(t *testing.T) {
		cf, err := GetConfig(&TestingConfigProvider{})
		if err != nil {
			t.Errorf("GetConfig returned error %s", err)
		}

		if cf.YoutubeAPIKey != ytTestKey {
			t.Errorf("YoutubeAPIKey is not the one provided by TestingConfigProvider. Expected %s got %s", ytTestKey, cf.YoutubeAPIKey)
		}
	})
}
