package main

import (
	"os"
	"testing"
	"time"
)

func TestGetVideoIDFromFileName(t *testing.T) {
	id := "OGK8gnP4TfA"
	test := "Test Video 1-OGK8gnP4TfA.mp4"
	testResult, err := getVideoIDFromFileName(test)
	if err != nil {
		t.Errorf("Video name returned error: %s", err)
	}
	if testResult != id {
		t.Errorf("Video name ID parser did not result in correct ID: Expected %s got %s", id, testResult)
	}

	id = "zVn7GctHoVQ"
	test = "Test Video - With a Dash-zVn7GctHoVQ.mp4"
	testResult, err = getVideoIDFromFileName(test)
	if err != nil {
		t.Errorf("Video name returned error: %s", err)
	}
	if testResult != id {
		t.Errorf("Video name with a dash did not result in correct ID: Expected %s got %s", id, testResult)
	}

	test = "Test Video - With no ID.mp4"
	testResult, err = getVideoIDFromFileName(test)
	if err == nil {
		t.Errorf("Video with no ID should return error, returned: %s", testResult)
	}

	id = "1kt7-O837H8"
	test = "Test Video - ID with dash 1kt7-O837H8.mp4"
	testResult, err = getVideoIDFromFileName(test)
	if err != nil {
		t.Errorf("Video name returned error: %s", err)
	}
	if testResult != id {
		t.Errorf("Video name with a dash in ID did not result in correct ID: Expected %s got %s", id, testResult)
	}

	id = "dmZR-LFp4ns"
	test = "Test Video - with an extra dot (0.8)-dmZR-LFp4ns.mp4"
	testResult, err = getVideoIDFromFileName(test)
	if err != nil {
		t.Errorf("Video name returned error: %s", err)
	}
	if testResult != id {
		t.Errorf("Video name with a dot in it did not result in correct ID: Expected %s got %s", id, testResult)
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
			"Bad File.description",
			31000000,
			false,
		},
		MockFileInfo{
			"The Macintosh Quadra 800--AC4HwzAK7A.mp4",
			31000000,
			false,
		},
		MockFileInfo{
			"Bad File.x",
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

func TestGetFileTypeValidType(t *testing.T) {
	file := "test.mp4"
	fileType, err := getFileType(file)
	if err != nil {
		t.Errorf("getFileType returned an error %s", err)
	}
	if fileType != "mp4" {
		t.Errorf("getFileType did not return correct result. Expecte mp4, got %s", fileType)
	}

	file = "test.MP4"
	fileType, err = getFileType(file)
	if err != nil {
		t.Errorf("getFileType returned an error %s", err)
	}
	if fileType != "mp4" {
		t.Errorf("getFileType did not return correct result. Expecte mp4, got %s", fileType)
	}
}

func TestGetFileTypeInvalidType(t *testing.T) {
	file := "test."
	_, err := getFileType(file)
	if err == nil {
		t.Errorf("getFileType did not return an error when provided filename %s", file)
	}

	file = "test"
	_, err = getFileType(file)
	if err == nil {
		t.Errorf("getFileType did not return an error when provided filename %s", file)
	}

	file = ""
	_, err = getFileType(file)
	if err == nil {
		t.Errorf("getFileType did not return an error when provided filename %s", file)
	}
}

func TestIsFileMP4(t *testing.T) {
	file := "test.mp4"
	isValid, err := isMP4(file)
	if err != nil {
		t.Errorf("isMP4 returned an error %s", err)
	}
	if !isValid {
		t.Errorf("isMP4 returned false for %s", file)
	}

	file = "test.doc"
	isValid, err = isMP4(file)
	if err != nil {
		t.Errorf("isMP4 returned an error %s", err)
	}
	if isValid {
		t.Errorf("isMP4 returned true for %s", file)
	}

	file = "test."
	_, err = isMP4(file)
	if err == nil {
		t.Error("isMP4 did not return an error")
	}

	file = "test"
	_, err = isMP4(file)
	if err == nil {
		t.Error("isMP4 did not return an error")
	}

	file = ""
	_, err = isMP4(file)
	if err == nil {
		t.Error("isMP4 did not return an error")
	}
}
