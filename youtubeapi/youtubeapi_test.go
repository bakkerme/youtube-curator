package youtubeapi

import (
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/config"
	"net/http"
	"reflect"
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

		url := getVideoInfoURL(&[]string{id}, accessKey)

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

var fakeErrorMessage = "The PuppyGirl did the loudest bark"

type mockHTTPClient struct {
	throwError          bool
	statusCodeToReturn  int
	malformJSONResponse bool
}

func (ht *mockHTTPClient) Get(url string) (*http.Response, []byte, error) {
	testFile, err := loadVideoListSingleTestFile()
	if err != nil {
		return nil, nil, err
	}

	if ht.throwError {
		return nil, nil, errors.New(fakeErrorMessage)
	}

	response := &http.Response{
		StatusCode: ht.statusCodeToReturn,
	}

	if ht.malformJSONResponse {
		testFile = []byte("{[sdfsd{")
	}

	return response, testFile, nil
}

func TestGetVideoInfo(t *testing.T) {
	t.Run("getVideoInfo returns valid video info", func(t *testing.T) {
		apiTestString := "123abc"

		config := config.Config{
			YoutubeAPIKey: apiTestString,
		}

		videoListResponse, err := getVideoInfo(&[]string{"18-elPdai_1"}, &config, &mockHTTPClient{statusCodeToReturn: 200})
		if err != nil {
			t.Errorf("Got error when running getVideoInfo %s", err)
		}

		if !reflect.DeepEqual(*videoListResponse, vlExpectedSingleVideo) {
			t.Errorf("VideoListResponse results are different.\nExpected %+v\nGot %+v", vlExpectedSingleVideo, videoListResponse)
		}
	})

	t.Run("getVideoInfo throws error if HTTP Get fails", func(t *testing.T) {
		apiTestString := "123abc"

		config := config.Config{
			YoutubeAPIKey: apiTestString,
		}

		_, err := getVideoInfo(&[]string{"18-elPdai_1"}, &config, &mockHTTPClient{throwError: true})
		if err == nil {
			t.Errorf("Did not recieve error when HTTP Request threw error")
		}

		if !strings.Contains(err.Error(), fakeErrorMessage) {
			t.Errorf("Did not recieve error correct error back. Got: %s, expected to see %s string", err, fakeErrorMessage)
		}
	})

	t.Run("getVideoInfo throws error if HTTP Get returns malformed JSON", func(t *testing.T) {
		apiTestString := "123abc"

		config := config.Config{
			YoutubeAPIKey: apiTestString,
		}

		_, err := getVideoInfo(&[]string{"18-elPdai_1"}, &config, &mockHTTPClient{malformJSONResponse: true})
		if err == nil {
			t.Errorf("Did not recieve error when HTTP Request returned malformed JSON")
		}
	})

	t.Run("getVideoInfo throws error if given more than 50 IDs", func(t *testing.T) {
		apiTestString := "123abc"

		config := config.Config{
			YoutubeAPIKey: apiTestString,
		}

		var ids []string
		for i := 0; i < 50; i++ {
			ids = append(ids, fmt.Sprint(i))
		}

		_, err := getVideoInfo(&ids, &config, &mockHTTPClient{throwError: false})
		if err == nil {
			t.Errorf("Did not recieve error when given more than 50 IDs")
		}

	})
}
