package youtubeapi

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestConvertVideoAPIResponse(t *testing.T) {
	t.Run("Correctly takes a JSON API response and parses it into expected VideoListResponse object", func(t *testing.T) {
		file, err := ioutil.ReadFile("./testfiles/videorequest.json")
		if err != nil {
			t.Errorf("Loading VideoRequest json failed: %s", err)
		}

		vl, err := convertVideoAPIResponse(string(file))
		if err != nil {
			t.Errorf("convertVideoAPIResponse returned an error %s", err)
		}

		if !reflect.DeepEqual(*vl, vlExpectedFull) {
			t.Errorf("VideoListResponse results are different.\nExpected %+v\ngot      %+v", vlExpectedFull, *vl)
		}
	})

	t.Run("Throws an error when invalid JSON API response is entered", func(t *testing.T) {
		_, err := convertVideoAPIResponse("")
		if err == nil {
			t.Error("convertVideoAPIResponse did not return an error with invalid JSON")
		}

		_, err = convertVideoAPIResponse("{{")
		if err == nil {
			t.Error("convertVideoAPIResponse did not return an error with invalid JSON")
		}
	})
}
