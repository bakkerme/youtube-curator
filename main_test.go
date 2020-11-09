package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
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

type MockFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (f MockFileInfo) Name() string {
	return f.name
}
func (f MockFileInfo) Size() int64 {
	return f.size
}
func (f MockFileInfo) Mode() os.FileMode {
	return 0
}
func (f MockFileInfo) ModTime() time.Time {
	return time.Now()
}
func (f MockFileInfo) IsDir() bool {
	return f.isDir
}
func (f MockFileInfo) Sys() interface{} {
	return nil
}

func TestGetLocalVideosFromDirListWithValidDirList(t *testing.T) {
	dirlist := []os.FileInfo{
		MockFileInfo{
			"The Macintosh LC-dCqJ6iPHus0.mp4",
			84000000,
			false,
		},
		MockFileInfo{
			"The Macintosh Quadra 800--AC4HwzAK7A.mp4",
			31000000,
			false,
		},
		MockFileInfo{
			"The Macintosh SE-_gPsIiKtybA.mp4",
			32000000,
			false,
		},
	}

	path := "/test/path"

	videoList := []Video{
		Video{
			path + "/" + dirlist[0].Name(),
			"dCqJ6iPHus0",
		},
		Video{
			path + "/" + dirlist[1].Name(),
			"-AC4HwzAK7A",
		},
		Video{
			path + "/" + dirlist[2].Name(),
			"_gPsIiKtybA",
		},
	}

	videos, err := getLocalVideosFromDirList(&dirlist, path)
	if err != nil {
		t.Errorf("getLocalVideosFromDirList returned an error %s", err)
	}

	if (*videos)[0] != videoList[0] &&
		(*videos)[1] != videoList[1] &&
		(*videos)[2] != videoList[2] {
		t.Errorf("Video list does not match response from function: %s", videos)
	}
}

func TestGetLocalVideosFromDirListWithInvalidDirList(t *testing.T) {
	dirlist := []os.FileInfo{}
	path := "/test/path"
	videos, err := getLocalVideosFromDirList(&dirlist, path)
	if err != nil {
		t.Errorf("getLocalVideosFromDirList returned an error %s", err)
	}

	if len(*videos) > 0 {
		t.Errorf("getLocalVideosFromDirList with no dir list returned something %s", videos)
	}
}

func TestConvertRSSStringToRSS(t *testing.T) {
	rssExpect := RSS{
		"yt:channel:UCS-WzPVpAAli-1IfEG2lN8A",
		"Test Guy",
		[]Entry{
			Entry{
				"yt:video:KQA9Na4aOa1",
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
			Entry{
				"yt:video:OGK8gnP4TfA",
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
			Entry{
				"yt:video:FazJqPQ6xSs",
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

	file, err := ioutil.ReadFile("./tests/test.xml")
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
}

func TestConvertRSSStringToRSSWithInvalidXML(t *testing.T) {
	_, err := convertRSSStringToRSS("")
	if err == nil {
		t.Error("convertRSSStringToRSS did not return an error with invalid XML")
	}
}
