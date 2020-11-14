package main

import (
	"io/ioutil"
	"reflect"
	"testing"
)

var rssExpect = RSS{
	"yt:channel:UCS-WzPVpAAli-1IfEG2lN8A",
	"Test Guy",
	[]VideoEntry{
		VideoEntry{
			"KQA9Na4aOa1",
			"Test Video New",
			Link{
				"https://www.youtube.com/watch?v=KQA9Na4aOa1",
			},
			"2020-11-06T19:00:01+00:00",
			"2020-11-06T23:12:15+00:00",
			MediaGroup{
				"Test Video 1",
				Thumbnail{
					"https://i2.ytimg.com/vi/KQA9Na4aOa1/hqdefault.jpg",
					480,
					360,
				},
				"Test Description New",
			},
		},
		VideoEntry{
			"OGK8gnP4TfA",
			"Test Video 1",
			Link{
				"https://www.youtube.com/watch?v=OGK8gnP4TfA",
			},
			"2020-11-06T19:00:01+00:00",
			"2020-11-06T23:12:15+00:00",
			MediaGroup{
				"Test Video 1",
				Thumbnail{
					"https://i2.ytimg.com/vi/OGK8gnP4TfA/hqdefault.jpg",
					480,
					360,
				},
				"Test Description",
			},
		},
		VideoEntry{
			"FazJqPQ6xSs",
			"Test Video 2",
			Link{
				"https://www.youtube.com/watch?v=FazJqPQ6xSs",
			},
			"2020-11-06T19:00:01+00:00",
			"2020-11-06T23:12:15+00:00",
			MediaGroup{
				"Test Video 2",
				Thumbnail{
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
