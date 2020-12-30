package config

import (
	"hyperfocus.systems/youtube-curator-server/testutils"
	"os"
	"reflect"
	"testing"
)

type TestingConfigProvider struct {
	missingValue bool
}

var ytTestKey string = "123abc"
var videoPath string = "/home/videos"

func (cp TestingConfigProvider) LoadConfig() (*Config, error) {
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
	})
}

func TestEnvarConfigProvider(t *testing.T) {
	t.Run("Returns a set key envar", func(t *testing.T) {
		key := "TEST_KEY"
		value := "VALUE"
		err := os.Setenv(key, value)
		if err != nil {
			t.Error(err)
		}

		ecp := &EnvarConfigProvider{}
		result, err := ecp.getValue(key)
		if err != nil {
			t.Error(err)
		}

		if result != value {
			t.Errorf("EnvarConfigProvider GetValue did not present a correct value. Expected %s, got %s", value, result)
		}
	})

	t.Run("Throws an error if value is not set", func(t *testing.T) {
		ecp := &EnvarConfigProvider{}
		_, err := ecp.getValue("SOME_NONEXISTENT_VALUE")
		if err == nil {
			t.Errorf("EnvarConfigProvider should have returned an error for a non-existent-value")
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
