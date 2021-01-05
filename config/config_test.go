package config

import (
	"errors"
	"hyperfocus.systems/youtube-curator-server/testutils"
	"hyperfocus.systems/youtube-curator-server/utils"
	"reflect"
	"testing"
)

var ytTestKey string = "123abc"
var videoPath string = "/home/videos/"

type TestingConfigProvider struct {
	returnError  bool
	returnConfig *Config
}

func (tcp TestingConfigProvider) LoadConfig() (*Config, error) {
	if tcp.returnError {
		return nil, errors.New("An error is here")
	}

	if tcp.returnConfig != nil {
		return tcp.returnConfig, nil
	}

	return &Config{
		YoutubeAPIKey: ytTestKey,
		VideoDirPath:  videoPath,
	}, nil
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

		if cf.VideoDirPath != videoPath {
			t.Errorf("VideoDirPath is not the one provided by TestingConfigProvider. Expected %s got %s", videoPath, cf.VideoDirPath)
		}
	})

	t.Run("Adds a / to the end of VideoDirPath without it", func(t *testing.T) {
		path := "/a/path"
		cf, err := GetConfig(&TestingConfigProvider{
			returnConfig: &Config{
				YoutubeAPIKey: "asdf123",
				VideoDirPath:  path,
			},
		})
		if err != nil {
			t.Errorf("GetConfig returned error %s", err)
		}

		if cf.VideoDirPath != path+"/" {
			t.Errorf("A / was not added to the end of the video path")
		}
	})

	t.Run("Returns an error if LoadConfig fails", func(t *testing.T) {
		_, err := GetConfig(&TestingConfigProvider{returnError: true})
		if err == nil {
			t.Errorf("Expected an error to be returned")
		}
	})

	t.Run("Returns an error if YoutubeAPIKey field is invalid", func(t *testing.T) {
		_, err := GetConfig(&TestingConfigProvider{
			returnConfig: &Config{
				YoutubeAPIKey: "",
				VideoDirPath:  "/a/test",
			},
		})
		if err == nil {
			t.Errorf("Expected an error to be returned")
		}
	})

	t.Run("Returns an error if VideoDirPath field is invalid", func(t *testing.T) {
		_, err := GetConfig(&TestingConfigProvider{
			returnConfig: &Config{
				YoutubeAPIKey: "123abc",
				VideoDirPath:  "",
			},
		})
		if err == nil {
			t.Errorf("Expected an error to be returned")
		}
	})

	t.Run("Returns an error if VideoDirPath field is invalid", func(t *testing.T) {
		_, err := GetConfig(&TestingConfigProvider{
			returnConfig: &Config{
				YoutubeAPIKey: "123abc",
				VideoDirPath:  "",
			},
		})
		if err == nil {
			t.Errorf("Expected an error to be returned")
		}
	})
}

func TestEnvarConfigProvider(t *testing.T) {
	t.Run("LoadConfig runs correctly", func(t *testing.T) {
		expectConfig := &Config{
			YoutubeAPIKey: "123abc",
			VideoDirPath:  "/a/path",
		}

		ecp := &EnvarConfigProvider{}
		cfg, err := ecp.loadConfig(&utils.MockEnvRead{
			ReturnValueForInput: map[string]string{
				"YOUTUBE_API_KEY": "123abc",
				"VIDEO_DIR_PATH":  "/a/path",
			},
		})
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(expectConfig, cfg) {
			t.Errorf("EnvarConfigProvider did not provide correct config. Expected\n%s\ngot\n%s", expectConfig, cfg)
		}
	})

	t.Run("LoadConfig returns an error when it can't find YOUTUBE_API_KEY", func(t *testing.T) {
		ecp := &EnvarConfigProvider{}
		_, err := ecp.loadConfig(&utils.MockEnvRead{
			ReturnValueForInput: map[string]string{
				"VIDEO_DIR_PATH": "/a/path",
			},
		})

		if err == nil {
			t.Error("loadConfig should have returned error")
		}
	})

	t.Run("LoadConfig returns an error when it can't find VIDEO_DIR_PATH", func(t *testing.T) {
		ecp := &EnvarConfigProvider{}
		_, err := ecp.loadConfig(&utils.MockEnvRead{
			ReturnValueForInput: map[string]string{
				"YOUTUBE_API_KEY": "123abc",
			},
		})

		if err == nil {
			t.Error("loadConfig should have returned error")
		}
	})
}

func TestFileConfigProvider(t *testing.T) {
	t.Run("Loads a valid config JSON", func(t *testing.T) {
		home := "/home/testhome"

		expectedFilename := home + "/.config/yt-up2date/config.json"
		fileJSON := []byte(`{
			"youtubeAPIKey": "1234567654",
			"videoDirPath": "/a/path"
		}`)

		expectedConfig := &Config{
			YoutubeAPIKey: "1234567654",
			VideoDirPath:  "/a/path",
		}

		dr := testutils.MockDirReader{
			T:                   t,
			ExpectedFilename:    &expectedFilename,
			ReturnReadFileValue: &fileJSON,
			ReturnHomeDirPath:   &home,
		}

		cp := FileConfigProvider{}
		cfg, err := cp.loadConfig(&dr)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(*cfg, *expectedConfig) {
			t.Errorf("FileConfigProvider did not return expected result. Expected\n%s\ngot\n%s", *expectedConfig, *cfg)
		}
	})

	t.Run("Errors when file cannot be found", func(t *testing.T) {
		home := "/home/testhome"
		dr := testutils.MockDirReader{
			T:                   t,
			ShouldErrorReadFile: true,
			ReturnHomeDirPath:   &home,
		}

		cp := FileConfigProvider{}
		cfg, err := cp.loadConfig(&dr)

		if err == nil {
			t.Errorf("FileConfigProvider should have returned error. Got result of %+v", cfg)
		}
	})

	t.Run("Errors on invalid JSON", func(t *testing.T) {
		home := "/home/testhome"
		fileJSON := []byte(`{`)

		dr := testutils.MockDirReader{
			T:                   t,
			ReturnReadFileValue: &fileJSON,
			ReturnHomeDirPath:   &home,
		}

		cp := FileConfigProvider{}
		cfg, err := cp.loadConfig(&dr)

		if err == nil {
			t.Errorf("FileConfigProvider should have returned error. Got result of %+v", cfg)
		}
	})
}
