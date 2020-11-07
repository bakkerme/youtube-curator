package main

import (
	"testing"
)

func TestGetVideoIDFromFileName(t *testing.T) {
	id1 := "OGK8gnP4TfA"
	test1 := "Test Video 1-OGK8gnP4TfA.mp4"

	testResult1, err := getVideoIDFromFileName(test1)
	if err != nil {
		t.Errorf("Video name returned error: %s", err)
	}
	if testResult1 != id1 {
		t.Errorf("Video name ID parser did not result in correct ID: Expected %s got %s", id1, testResult1)
	}

	id2 := "zVn7GctHoVQ"
	test2 := "Test Video - With a Dash-zVn7GctHoVQ.mp4"
	testResult2, err := getVideoIDFromFileName(test2)
	if err != nil {
		t.Errorf("Video name returned error: %s", err)
	}
	if testResult2 != id2 {
		t.Errorf("Video name with a dash did not result in correct ID: Expected %s got %s", id1, testResult1)
	}

	test3 := "Test Video - With no ID.mp4"
	testResult3, err := getVideoIDFromFileName(test3)
	if err == nil {
		t.Errorf("Video with no ID should return error, returned: %s", testResult3)
	}

	id3 := "1kt7-O837H8"
	test4 := "Test Video - ID with dash 1kt7-O837H8.mp4"
	testResult4, err := getVideoIDFromFileName(test4)
	if testResult4 != id3 {
		t.Errorf("Video name with a dash in ID did not result in correct ID: Expected %s got %s", id3, testResult4)
	}
}

func TestGetEntriesNotInVideoListWithMatch(t *testing.T) {
	thumbnail := Thumbnail{"", 0, 0}
	mediaGroup := MediaGroup{"", thumbnail, ""}

	outstandingEntry := Entry{"wxyzabcdefg", "Video 4", Link{"http://link4"}, "", "", mediaGroup}
	entries := []Entry{
		Entry{"12345678911", "Video 1", Link{"http://link"}, "", "", mediaGroup},
		Entry{"abcdefghijk", "Video 2", Link{"http://link2"}, "", "", mediaGroup},
		Entry{"lmnopqrxtuv", "Video 3", Link{"http://link3"}, "", "", mediaGroup},
		outstandingEntry,
	}

	videos := []Video{
		Video{"Video 1-12345678911.mp4", "12345678911"},
		Video{"Video 2-abcdefghijk.mp4", "abcdefghijk"},
		Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv"},
	}

	entriesToGet := getEntriesNotInVideoList(&entries, &videos)

	if (*entriesToGet)[0] != outstandingEntry {
		t.Errorf("The outstanding entry not in the video list is incorrect, got %s, expected %s", (*entriesToGet)[0].ID, outstandingEntry.ID)
	}
}

func TestGetEntriesNotInVideoListWithNoMatch(t *testing.T) {
	thumbnail := Thumbnail{"", 0, 0}
	mediaGroup := MediaGroup{"", thumbnail, ""}

	entries := []Entry{
		Entry{"12345678911", "Video 1", Link{"http://link"}, "", "", mediaGroup},
		Entry{"abcdefghijk", "Video 2", Link{"http://link2"}, "", "", mediaGroup},
		Entry{"lmnopqrxtuv", "Video 3", Link{"http://link3"}, "", "", mediaGroup},
	}

	videos := []Video{
		Video{"Video 1-12345678911.mp4", "12345678911"},
		Video{"Video 2-abcdefghijk.mp4", "abcdefghijk"},
		Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv"},
	}

	entriesToGet := getEntriesNotInVideoList(&entries, &videos)

	if len(*entriesToGet) > 0 {
		t.Errorf("Unknown match was found, got %s", (*entriesToGet)[0].ID)
	}
}

func TestGetEntriesNotInVideoListWithNoVideoList(t *testing.T) {
	thumbnail := Thumbnail{"", 0, 0}
	mediaGroup := MediaGroup{"", thumbnail, ""}

	entries := []Entry{
		Entry{"12345678911", "Video 1", Link{"http://link"}, "", "", mediaGroup},
		Entry{"abcdefghijk", "Video 2", Link{"http://link2"}, "", "", mediaGroup},
		Entry{"lmnopqrxtuv", "Video 3", Link{"http://link3"}, "", "", mediaGroup},
	}

	videos := []Video{}

	entriesToGet := getEntriesNotInVideoList(&entries, &videos)

	if len(*entriesToGet) != 3 {
		t.Errorf("Matches were not correctly found")
	}
}

func TestGetEntriesNotInVideoListWithNoEntries(t *testing.T) {
	entries := []Entry{}

	videos := []Video{
		Video{"Video 1-12345678911.mp4", "12345678911"},
		Video{"Video 2-abcdefghijk.mp4", "abcdefghijk"},
		Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv"},
	}

	entriesToGet := getEntriesNotInVideoList(&entries, &videos)

	if len(*entriesToGet) > 0 {
		t.Errorf("Unknown match was found, got %s", (*entriesToGet)[0].ID)
	}
}

func TestIsEntryInVideoListWithMatch(t *testing.T) {
	thumbnail := Thumbnail{"", 0, 0}
	mediaGroup := MediaGroup{"", thumbnail, ""}
	entry := Entry{"12345678911", "Video 1", Link{"http://link"}, "", "", mediaGroup}

	videos := []Video{
		Video{"Video 1-12345678911.mp4", "12345678911"},
		Video{"Video 2-abcdefghijk.mp4", "abcdefghijk"},
		Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv"},
	}

	if !isEntryInVideoList(&entry, &videos) {
		t.Error("Entry should have been found in video list")
	}
}

func TestIsEntryInVideoListWithNoMatch(t *testing.T) {
	thumbnail := Thumbnail{"", 0, 0}
	mediaGroup := MediaGroup{"", thumbnail, ""}
	entry := Entry{"BADID123456", "Video 1", Link{"http://link"}, "", "", mediaGroup}

	videos := []Video{
		Video{"Video 1-12345678911.mp4", "12345678911"},
		Video{"Video 2-abcdefghijk.mp4", "abcdefghijk"},
		Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv"},
	}

	if isEntryInVideoList(&entry, &videos) {
		t.Error("Entry should not have been found in video list")
	}
}

func TestIsEntryInVideoListWithEmptyVideoList(t *testing.T) {
	thumbnail := Thumbnail{"", 0, 0}
	mediaGroup := MediaGroup{"", thumbnail, ""}
	entry := Entry{"BADID123456", "Video 1", Link{"http://link"}, "", "", mediaGroup}

	videos := []Video{}

	if isEntryInVideoList(&entry, &videos) {
		t.Error("Entry should not have been found in video list")
	}
}
