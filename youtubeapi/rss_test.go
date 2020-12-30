package youtubeapi

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

var rssExpect = RSS{
	"yt:channel:UCS-WzPVpAAli-1IfEG2lN8A",
	"Test Guy",
	[]RSSVideoEntry{
		RSSVideoEntry{
			"KQA9Na4aOa1",
			"Test Video New",
			RSSLink{
				"https://www.youtube.com/watch?v=KQA9Na4aOa1",
			},
			"2020-11-06T19:00:01+00:00",
			"2020-11-06T23:12:15+00:00",
			RSSMediaGroup{
				"Test Video 1",
				RSSThumbnail{
					"https://i2.ytimg.com/vi/KQA9Na4aOa1/hqdefault.jpg",
					480,
					360,
				},
				"Test Description New",
			},
		},
		RSSVideoEntry{
			"OGK8gnP4TfA",
			"Test Video 1",
			RSSLink{
				"https://www.youtube.com/watch?v=OGK8gnP4TfA",
			},
			"2020-11-06T19:00:01+00:00",
			"2020-11-06T23:12:15+00:00",
			RSSMediaGroup{
				"Test Video 1",
				RSSThumbnail{
					"https://i2.ytimg.com/vi/OGK8gnP4TfA/hqdefault.jpg",
					480,
					360,
				},
				"Test Description",
			},
		},
		RSSVideoEntry{
			"FazJqPQ6xSs",
			"Test Video 2",
			RSSLink{
				"https://www.youtube.com/watch?v=FazJqPQ6xSs",
			},
			"2020-11-06T19:00:01+00:00",
			"2020-11-06T23:12:15+00:00",
			RSSMediaGroup{
				"Test Video 2",
				RSSThumbnail{
					"https://i2.ytimg.com/vi/FazJqPQ6xSs/hqdefault.jpg",
					480,
					360,
				},
				"Test Description 2",
			},
		},
	},
}

func TestConvertRSSStringToRSS(t *testing.T) {
	t.Run("Correctly takes RSS feed and parses into expected RSS object", func(t *testing.T) {
		file, err := ioutil.ReadFile("./testfiles/test.xml")
		if err != nil {
			t.Errorf("Loading RSS feed xml failed: %s", err)
		}

		rss, err := convertRSSStringToRSS(string(file))
		if err != nil {
			t.Errorf("convertRSSStringToRSS returned an error %s", err)
		}

		if !reflect.DeepEqual(*rss, rssExpect) {
			t.Errorf("RSS results are different. Expected %+v, got %+v", rssExpect, *rss)
		}
	})

	t.Run("Throws an error when invalid RSS is entered", func(t *testing.T) {
		_, err := convertRSSStringToRSS("")
		if err == nil {
			t.Error("convertRSSStringToRSS did not return an error with invalid XML")
		}

		_, err = convertRSSStringToRSS("<aaaa")
		if err == nil {
			t.Error("convertRSSStringToRSS did not return an error with invalid XML")
		}
	})
}

func TestGetRSSFeed(t *testing.T) {
	t.Run("GetRSSFeed pulls data over HTTP and converts it to RSS", func(t *testing.T) {
		rss, err := GetRSSFeed("testurl.homebase", &mockHTTPClient{statusCodeToReturn: 200, responseFile: "./testfiles/test.xml"})

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(*rss, rssExpect) {
			t.Errorf("RSS results are different. Expected %+v, got %+v", rssExpect, *rss)
		}
	})

	t.Run("GetRSSFeed returns an error if HTTP Client returns an error", func(t *testing.T) {
		_, err := GetRSSFeed("testurl.homebase", &mockHTTPClient{throwError: true})

		if err == nil {
			t.Errorf("Expected GetRSSFeed would return an error")
		}
	})

	t.Run("GetRSSFeed returns an error if HTTP status code is not 200", func(t *testing.T) {
		_, err := GetRSSFeed("testurl.homebase", &mockHTTPClient{statusCodeToReturn: 400})

		fmt.Println(err)
		if err == nil {
			t.Errorf("Expected GetRSSFeed would return an error")
		}
	})
}
