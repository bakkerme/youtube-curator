package config

import (
	"fmt"
	"os"
)

// Config represents application-level configuration
type Config struct {
	YoutubeAPIKey string
	VideoDirPath  string
}

// configProvider is an interface for providers of the configuration
type configProvider interface {
	GetValue(string) (string, error)
}

// GetConfig returns a Config struct containing application-level configuration
func GetConfig(cp configProvider) (*Config, error) {
	youtubeAPIKey, err := cp.GetValue("YOUTUBE_API_KEY")
	if err != nil {
		return nil, err
	}
	dirPath, err := cp.GetValue("VIDEO_DIR_PATH")
	if err != nil {
		return nil, err
	}

	return &Config{
		YoutubeAPIKey: youtubeAPIKey,
		VideoDirPath:  dirPath,
	}, nil
}

// EnvarConfigProvider provides configuration from the environment variables
type EnvarConfigProvider struct{}

func (cp EnvarConfigProvider) GetValue(key string) (string, error) {
	result, didFind := os.LookupEnv(key)
	if didFind {
		return result, nil
	}

	return "", fmt.Errorf("Could not find %s in Environment for Config", key)
}
