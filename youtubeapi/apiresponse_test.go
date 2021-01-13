package youtubeapi

import (
	"hyperfocus.systems/youtube-curator-server/testutils"
	"reflect"
	"testing"
)

func TestConvertVideoAPIResponse(t *testing.T) {
	t.Run("Correctly takes a JSON API response and parses it into expected VideoListResponse object", func(t *testing.T) {
		vl, err := convertAPIResponse(string(SearchResponseJSON), apiSearch)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("convertAPIResponse", err))
		}

		if !reflect.DeepEqual(*vl, expectedSearchResponse) {
			t.Errorf(testutils.MismatchError("convertAPIResponse", expectedSearchResponse, *vl))
		}
	})

	t.Run("Throws an error when invalid JSON API response is entered", func(t *testing.T) {
		_, err := convertAPIResponse("", apiSearch)
		if err == nil {
			t.Errorf(testutils.ExpectedError("convertAPIResponse"))
		}

		_, err = convertAPIResponse("{{", apiSearch)
		if err == nil {
			t.Errorf(testutils.ExpectedError("convertAPIResponse"))
		}
	})

	t.Run("converts Video JSON correctly", func(t *testing.T) {
		vl, err := convertAPIResponse(string(VideoResponseJSON), apiVideos)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("convertAPIResponse", err))
		}

		if !reflect.DeepEqual(*vl, expectedVideoResponse) {
			t.Errorf(testutils.MismatchError("convertAPIResponse", expectedSearchResponse, *vl))
		}
	})

	t.Run("converts PlaylistItems JSON correctly", func(t *testing.T) {
		vl, err := convertAPIResponse(string(PlaylistItemsResponseJSON), apiPlaylistItems)
		if err != nil {
			t.Errorf(testutils.UnexpectedError("convertAPIResponse", err))
		}

		if !reflect.DeepEqual(*vl, expectedPlaylistItemsResponse) {
			t.Errorf(testutils.MismatchError("convertAPIResponse", expectedPlaylistItemsResponse, *vl))
		}
	})
}
