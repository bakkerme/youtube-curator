package config

import (
	"fmt"
	"os"
	"testing"
)

type TestingConfigProvider struct {
	missingValue bool
}

var ytTestKey string = "123abc"

func (cp TestingConfigProvider) getValue(key string) (string, error) {
	if cp.missingValue {
		return "", fmt.Errorf("Correctly missing value")
	}

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

	t.Run("Returns an error if a value is not available", func(t *testing.T) {
		_, err := GetConfig(&TestingConfigProvider{missingValue: true})
		if err == nil {
			t.Errorf("GetConfig should have returned an error")
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
			t.Errorf("EnvarConfigProvider getValue did not present a correct value. Expected %s, got %s", value, result)
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
