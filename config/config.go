package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/utils"
)

// Config represents application-level configuration
type Config struct {
	YoutubeAPIKey string `json:"youtubeAPIKey"`
	VideoDirPath  string `json:"videoDirPath"`
}

// configProvider is an interface for providers of the configuration
type configProvider interface {
	LoadConfig() (*Config, error)
}

// GetConfig returns a Config struct containing application-level configuration
func GetConfig(cp configProvider) (*Config, error) {
	cfg, err := cp.LoadConfig()
	if err != nil {
		return nil, err
	}

	err = checkFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("Config format is invalid.\n%s", err)
	}

	dirPath := cfg.VideoDirPath
	if dirPath[len(dirPath)-1] != '/' {
		cfg.VideoDirPath = dirPath + "/"
	}

	return cfg, nil
}

func checkFields(cfg *Config) error {
	if len(cfg.YoutubeAPIKey) == 0 {
		return errors.New("YoutubeAPIKey in config is invalid. It should be a 39 character long string. See readme for more info")
	}

	if len(cfg.VideoDirPath) == 0 {
		return errors.New("VideoDirPath in config is invalid. This should be a path to a folder intended to store videos")
	}

	return nil
}

// EnvarConfigProvider provides configuration from the environment variables
type EnvarConfigProvider struct{}

// LoadConfig loads a config file from the set environment variables
func (cp EnvarConfigProvider) LoadConfig() (*Config, error) {
	return cp.loadConfig(&utils.EnvRead{})
}

func (cp EnvarConfigProvider) loadConfig(envr utils.EnvReader) (*Config, error) {
	youtubeAPIKey, didFind := envr.LookupEnv("YOUTUBE_API_KEY")
	if !didFind {
		return nil, errors.New("Could not find YOUTUBE_API_KEY")
	}

	videoDirPath, didFind := envr.LookupEnv("VIDEO_DIR_PATH")
	if !didFind {
		return nil, errors.New("Could not find VIDEO_DIR_PATH")
	}

	return &Config{
		YoutubeAPIKey: youtubeAPIKey,
		VideoDirPath:  videoDirPath,
	}, nil
}

// FileConfigProvider provides configuration from the environment variables
type FileConfigProvider struct {
	config *Config
}

// LoadConfig will load a config file off disk
func (cp FileConfigProvider) LoadConfig() (*Config, error) {
	return cp.loadConfig(&utils.DirReader{})
}

func (cp FileConfigProvider) loadConfig(dr utils.DirReaderProvider) (*Config, error) {
	path := dr.GetHomeDirPath() + "/.config/yt-up2date/config.json"
	file, err := dr.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Could not find Config file, please see readme and add valid config.json to %s. Error %s", path, err)
	}

	var cfg Config
	if err := json.Unmarshal([]byte(file), &cfg); err != nil {
		return nil, fmt.Errorf("Can't unmarshal config file. Loading up %s, got error %s", path, err)
	}

	return &cfg, nil
}
