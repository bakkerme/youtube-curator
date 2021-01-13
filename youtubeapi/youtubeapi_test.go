package youtubeapi

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/testutils"
	"hyperfocus.systems/youtube-curator-server/utils"
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

func TestGetVideoMetadata(t *testing.T) {
	t.Run("getVideoMetadata returns valid video info", func(t *testing.T) {
		apiTestString := "123abc"

		config := config.Config{
			YoutubeAPIKey: apiTestString,
		}

		videoListResponse, err := getVideoMetadata(
			&[]string{"18-elPdai_1"},
			&config,
			&utils.MockHTTPClient{
				StatusCode: 200,
				Body:       []byte(VideoResponseJSON),
			})
		if err != nil {
			t.Errorf("Got error when running getVideoMetadata %s", err)
		}

		if !reflect.DeepEqual(*videoListResponse, expectedVideoResponse) {
			t.Errorf("VideoListResponse results are different.\nExpected %+v\nGot %+v", expectedVideoResponse, videoListResponse)
		}
	})

	t.Run("getVideoMetadata throws error if HTTP Get fails", func(t *testing.T) {
		apiTestString := "123abc"

		config := config.Config{
			YoutubeAPIKey: apiTestString,
		}

		_, err := getVideoMetadata(
			&[]string{"18-elPdai_1"},
			&config,
			&utils.MockHTTPClient{
				ThrowError: true,
				Body:       []byte(VideoResponseJSON),
			},
		)
		if err == nil {
			t.Errorf("Did not recieve error when HTTP Request threw error")
		}
	})

	t.Run("getVideoMetadata throws error if HTTP Get returns malformed JSON", func(t *testing.T) {
		apiTestString := "123abc"

		config := config.Config{
			YoutubeAPIKey: apiTestString,
		}

		_, err := getVideoMetadata(
			&[]string{"18-elPdai_1"},
			&config,
			&utils.MockHTTPClient{
				Body: []byte("{sdfds[{{"),
			},
		)
		if err == nil {
			t.Errorf("Did not recieve error when HTTP Request returned malformed JSON")
		}
	})

	t.Run("getVideoMetadata throws error if HTTP Get returns 400 status code", func(t *testing.T) {
		apiTestString := "123abc"

		config := config.Config{
			YoutubeAPIKey: apiTestString,
		}

		_, err := getVideoMetadata(
			&[]string{"18-elPdai_1"},
			&config,
			&utils.MockHTTPClient{
				StatusCode: 400,
			},
		)
		if err == nil {
			t.Errorf("Did not recieve error when HTTP Request returned 400")
		}
	})

	t.Run("getVideoMetadata throws error if given more than 50 IDs", func(t *testing.T) {
		apiTestString := "123abc"

		config := config.Config{
			YoutubeAPIKey: apiTestString,
		}

		var ids []string
		for i := 0; i < 51; i++ {
			ids = append(ids, fmt.Sprint(i))
		}

		_, err := getVideoMetadata(
			&ids,
			&config,
			&utils.MockHTTPClient{
				ThrowError: false,
				Body:       []byte(VideoResponseJSON),
			},
		)
		if err == nil {
			t.Errorf("Did not recieve error when given more than 50 IDs")
		}
	})
}

func TestGetVideosForChannel(t *testing.T) {
	ytc := collection.MockYTChannel{
		IName:         "Name",
		IID:           "ID",
		IRSSURL:       "http://rssurl",
		IChannelURL:   "http://channelurl",
		IArchivalMode: collection.ArchivalModeArchive,
		IChannelType:  collection.ChannelTypeChannel,
	}

	cf := config.Config{
		YoutubeAPIKey: "ASDF123",
	}

	t.Run("getVideosForChannel returns valid results", func(t *testing.T) {
		httpClient := utils.MockHTTPClient{
			StatusCode: 200,
			Body:       []byte(SearchResponseJSON),
		}

		videoListResponse, err := getVideosForChannel(ytc.ChannelType(), &ytc, &cf, &httpClient)

		if err != nil {
			t.Errorf(testutils.UnexpectedError("getVideosForChannel", err))
		}

		if !reflect.DeepEqual(*videoListResponse, expectedSearchResponse) {
			t.Errorf(testutils.MismatchError("getVideosForChannel", expectedSearchResponse, videoListResponse))
		}
	})

	t.Run("getVideosForChannel uses the search API for Channel Type Channel", func(t *testing.T) {
		validationFunction := func(url string) {
			if !strings.Contains(url, "search") {
				t.Errorf("URL %s does not contain search", url)
			}
		}

		httpClient := utils.MockHTTPClient{
			StatusCode: 200,
			Body:       []byte(SearchResponseJSON),
			Validate:   &validationFunction,
		}

		_, err := getVideosForChannel(ytc.ChannelType(), &ytc, &cf, &httpClient)

		if err != nil {
			t.Errorf("getVideosForChannel returned an unexpected error %s", err)
		}
	})

	t.Run("getVideosForChannel uses the search API for Channel Type Playlist", func(t *testing.T) {
		ytcPlaylist := collection.MockYTChannel{
			IName:         "Name",
			IID:           "ID",
			IRSSURL:       "http://rssurl",
			IChannelURL:   "http://channelurl",
			IArchivalMode: collection.ArchivalModeArchive,
			IChannelType:  collection.ChannelTypePlaylist,
		}

		validationFunction := func(url string) {
			if !strings.Contains(url, "playlist") {
				t.Errorf("URL %s does not contain playlist", url)
			}
		}

		httpClient := utils.MockHTTPClient{
			StatusCode: 200,
			Body:       []byte(PlaylistItemsResponseJSON),
			Validate:   &validationFunction,
		}

		_, err := getVideosForChannel(ytcPlaylist.ChannelType(), &ytcPlaylist, &cf, &httpClient)

		if err != nil {
			t.Errorf("getVideosForChannel returned an unexpected error %s", err)
		}
	})

	t.Run("getVideosForChannel returns error if HTTP response is not 200", func(t *testing.T) {
		httpClient := utils.MockHTTPClient{
			StatusCode: 500,
			Body:       []byte(""),
		}

		_, err := getVideosForChannel(ytc.ChannelType(), &ytc, &cf, &httpClient)

		if err == nil {
			t.Error("getVideosForChannel did not return expected error")
		}
	})

	t.Run("getVideosForChannel returns error if HTTP response is invalid", func(t *testing.T) {
		httpClient := utils.MockHTTPClient{
			StatusCode: 200,
			Body:       []byte("234sdfsadf"),
		}

		_, err := getVideosForChannel(ytc.ChannelType(), &ytc, &cf, &httpClient)

		if err == nil {
			t.Error("getVideosForChannel did not return expected error")
		}
	})

	t.Run("getVideosForChannel returns error if HTTP response is body is empty", func(t *testing.T) {
		httpClient := utils.MockHTTPClient{
			StatusCode: 200,
			Body:       []byte(""),
		}

		_, err := getVideosForChannel(ytc.ChannelType(), &ytc, &cf, &httpClient)

		if err == nil {
			t.Error("getVideosForChannel did not return expected error")
		}
	})
}

func TestMakeAPIRequest(t *testing.T) {
	t.Run("makeAPIRequest returns valid results", func(t *testing.T) {
		expectedBody := []byte("SOME body")

		keyVals := map[string]string{
			"param":  "value1",
			"param2": "value2",
		}

		expectedURLParts := []string{
			baseURL,
			"someapi?",
			"key=123abc&",
			"param=value1&",
			"param2=value2",
		}

		validationFunction := func(url string) {
			for _, part := range expectedURLParts {
				if !strings.Contains(url, part) {
					t.Errorf("URL %s does not contain an expected part %s", url, part)
				}
			}
		}

		mockHTTPClient := &utils.MockHTTPClient{
			Body:     expectedBody,
			Validate: &validationFunction,
		}

		body, err := makeAPIRequest("someapi", &keyVals, "123abc", mockHTTPClient)

		if err != nil {
			t.Errorf("makeAPIRequest returned unexpected error %s", err)
		}

		if string(expectedBody) != string(body) {
			t.Errorf("makeAPIRequest did not return correct result. Expected %s, got %s", expectedBody, body)
		}
	})

	t.Run("makeAPIRequest makes request with nil params", func(t *testing.T) {
		expectedBody := []byte("SOME body")

		expectedURL := baseURL + "someapi?key=123abc"

		validationFunction := func(url string) {
			if url != expectedURL {
				t.Errorf("Value passed to URL is invalid. Expected\n%s got\n%s", expectedURL, url)
			}
		}

		mockHTTPClient := &utils.MockHTTPClient{
			Body:     expectedBody,
			Validate: &validationFunction,
		}

		body, err := makeAPIRequest("someapi", nil, "123abc", mockHTTPClient)

		if err != nil {
			t.Errorf("makeAPIRequest returned unexpected error %s", err)
		}

		if string(expectedBody) != string(body) {
			t.Errorf("makeAPIRequest did not return correct result. Expected %s, got %s", expectedBody, body)
		}
	})

	t.Run("makeAPIRequest returns error if http client returns error", func(t *testing.T) {
		mockHTTPClient := &utils.MockHTTPClient{
			Body:       []byte(""),
			ThrowError: true,
		}

		_, err := makeAPIRequest("someapi", nil, "123abc", mockHTTPClient)

		if err == nil {
			t.Errorf("makeAPIRequest should have returned an error")
		}
	})

	t.Run("makeAPIRequest returns error if status code is not 200", func(t *testing.T) {
		mockHTTPClient := &utils.MockHTTPClient{
			Body:       []byte(""),
			StatusCode: 500,
		}

		_, err := makeAPIRequest("someapi", nil, "123abc", mockHTTPClient)

		if err == nil {
			t.Errorf("makeAPIRequest should have returned an error")
		}
	})

	t.Run("makeAPIRequest returns error if access key is empty", func(t *testing.T) {
		mockHTTPClient := &utils.MockHTTPClient{
			Body:       []byte(""),
			StatusCode: 500,
		}

		_, err := makeAPIRequest("someapi", nil, "", mockHTTPClient)

		if err == nil {
			t.Errorf("makeAPIRequest should have returned an error")
		}
	})
}
