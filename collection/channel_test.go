package collection

import (
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGetVideoIDFromFileName(t *testing.T) {
	t.Run("Parses ID from standard video title", func(t *testing.T) {
		id := "OGK8gnP4TfA"
		test := "Test Video 1-OGK8gnP4TfA.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err != nil {
			t.Errorf("Video name returned error: %s", err)
		}
		if testResult != id {
			t.Errorf("Video name ID parser did not result in correct ID: Expected %s got %s", id, testResult)
		}
	})

	t.Run("Parses ID from video title with a dash", func(t *testing.T) {
		id := "zVn7GctHoVQ"
		test := "Test Video - With a Dash-zVn7GctHoVQ.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err != nil {
			t.Errorf("Video name returned error: %s", err)
		}
		if testResult != id {
			t.Errorf("Video name with a dash did not result in correct ID: Expected %s got %s", id, testResult)
		}
	})

	t.Run("Throws error for video title with no ID", func(t *testing.T) {
		test := "Test Video - With no ID.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err == nil {
			t.Errorf("Video with no ID should return error, returned: %s", testResult)
		}
	})

	t.Run("Parses ID from video title with dash in ID", func(t *testing.T) {
		id := "1kt7-O837H8"
		test := "Test Video - ID with dash 1kt7-O837H8.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err != nil {
			t.Errorf("Video name returned error: %s", err)
		}
		if testResult != id {
			t.Errorf("Video name with a dash in ID did not result in correct ID: Expected %s got %s", id, testResult)
		}
	})

	t.Run("Parses ID from video title a dot in the title", func(t *testing.T) {
		id := "dmZR-LFp4ns"
		test := "Test Video - with an extra dot (0.8)-dmZR-LFp4ns.mp4"
		testResult, err := getVideoIDFromFileName(test)
		if err != nil {
			t.Errorf("Video name returned error: %s", err)
		}
		if testResult != id {
			t.Errorf("Video name with a dot in it did not result in correct ID: Expected %s got %s", id, testResult)
		}
	})
}

func TestGetEntriesNotInVideoList(t *testing.T) {
	t.Run("With a Match", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{"", 0, 0}
		mediaGroup := youtubeapi.RSSMediaGroup{"", thumbnail, ""}

		outstandingEntry := youtubeapi.RSSVideoEntry{"wxyzabcdefg", "Video 4", youtubeapi.RSSLink{"http://link4"}, "", "", mediaGroup}
		entries := []youtubeapi.RSSVideoEntry{
			youtubeapi.RSSVideoEntry{"12345678911", "Video 1", youtubeapi.RSSLink{"http://link"}, "", "", mediaGroup},
			youtubeapi.RSSVideoEntry{"abcdefghijk", "Video 2", youtubeapi.RSSLink{"http://link2"}, "", "", mediaGroup},
			youtubeapi.RSSVideoEntry{"lmnopqrxtuv", "Video 3", youtubeapi.RSSLink{"http://link3"}, "", "", mediaGroup},
			outstandingEntry,
		}

		videos := []Video{
			Video{"Video 1-12345678911.mp4", "12345678911", "mp4", "/"},
			Video{"Video 2-abcdefghijk.mp4", "abcdefghijk", "mp4", "/"},
			Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv", "mp4", "/"},
		}

		entriesToGet := GetEntriesNotInVideoList(&entries, &videos)

		if (*entriesToGet)[0] != outstandingEntry {
			t.Errorf("The outstanding entry not in the video list is incorrect, got %s, expected %s", (*entriesToGet)[0].ID, outstandingEntry.ID)
		}
	})

	t.Run("With no match", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{"", 0, 0}
		mediaGroup := youtubeapi.RSSMediaGroup{"", thumbnail, ""}

		entries := []youtubeapi.RSSVideoEntry{
			youtubeapi.RSSVideoEntry{"12345678911", "Video 1", youtubeapi.RSSLink{"http://link"}, "", "", mediaGroup},
			youtubeapi.RSSVideoEntry{"abcdefghijk", "Video 2", youtubeapi.RSSLink{"http://link2"}, "", "", mediaGroup},
			youtubeapi.RSSVideoEntry{"lmnopqrxtuv", "Video 3", youtubeapi.RSSLink{"http://link3"}, "", "", mediaGroup},
		}

		videos := []Video{
			Video{"Video 1-12345678911.mp4", "12345678911", "mp4", "/"},
			Video{"Video 2-abcdefghijk.mp4", "abcdefghijk", "mp4", "/"},
			Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv", "mp4", "/"},
		}

		entriesToGet := GetEntriesNotInVideoList(&entries, &videos)

		if len(*entriesToGet) > 0 {
			t.Errorf("Unknown match was found, got %s", (*entriesToGet)[0].ID)
		}
	})

	t.Run("With no video list provided", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{"", 0, 0}
		mediaGroup := youtubeapi.RSSMediaGroup{"", thumbnail, ""}

		entries := []youtubeapi.RSSVideoEntry{
			youtubeapi.RSSVideoEntry{"12345678911", "Video 1", youtubeapi.RSSLink{"http://link"}, "", "", mediaGroup},
			youtubeapi.RSSVideoEntry{"abcdefghijk", "Video 2", youtubeapi.RSSLink{"http://link2"}, "", "", mediaGroup},
			youtubeapi.RSSVideoEntry{"lmnopqrxtuv", "Video 3", youtubeapi.RSSLink{"http://link3"}, "", "", mediaGroup},
		}

		videos := []Video{}

		entriesToGet := GetEntriesNotInVideoList(&entries, &videos)

		if len(*entriesToGet) != 3 {
			t.Errorf("Matches were not correctly found")
		}
	})

	t.Run("With no entries provied", func(t *testing.T) {
		entries := []youtubeapi.RSSVideoEntry{}

		videos := []Video{
			Video{"Video 1-12345678911.mp4", "12345678911", "mp4", "/"},
			Video{"Video 2-abcdefghijk.mp4", "abcdefghijk", "mp4", "/"},
			Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv", "mp4", "/"},
		}

		entriesToGet := GetEntriesNotInVideoList(&entries, &videos)

		if len(*entriesToGet) > 0 {
			t.Errorf("Unknown match was found, got %s", (*entriesToGet)[0].ID)
		}
	})

}

func TestIsEntryInVideoList(t *testing.T) {
	t.Run("With match", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{"", 0, 0}
		mediaGroup := youtubeapi.RSSMediaGroup{"", thumbnail, ""}
		entry := youtubeapi.RSSVideoEntry{"12345678911", "Video 1", youtubeapi.RSSLink{"http://link"}, "", "", mediaGroup}

		videos := []Video{
			Video{"Video 1-12345678911.mp4", "12345678911", "mp4", "/"},
			Video{"Video 2-abcdefghijk.mp4", "abcdefghijk", "mp4", "/"},
			Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv", "mp4", "/"},
		}

		if !isEntryInVideoList(&entry, &videos) {
			t.Error("youtubeapi.RSSVideoEntry should have been found in video list")
		}
	})

	t.Run("With no match", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{"", 0, 0}
		mediaGroup := youtubeapi.RSSMediaGroup{"", thumbnail, ""}
		entry := youtubeapi.RSSVideoEntry{"BADID123456", "Video 1", youtubeapi.RSSLink{"http://link"}, "", "", mediaGroup}

		videos := []Video{
			Video{"Video 1-12345678911.mp4", "12345678911", "mp4", "/"},
			Video{"Video 2-abcdefghijk.mp4", "abcdefghijk", "mp4", "/"},
			Video{"Video 3-lmnopqrxtuv.mp4", "lmnopqrxtuv", "mp4", "/"},
		}

		if isEntryInVideoList(&entry, &videos) {
			t.Error("youtubeapi.RSSVideoEntry should not have been found in video list")
		}
	})

	t.Run("With empty video list", func(t *testing.T) {
		thumbnail := youtubeapi.RSSThumbnail{"", 0, 0}
		mediaGroup := youtubeapi.RSSMediaGroup{"", thumbnail, ""}
		entry := youtubeapi.RSSVideoEntry{"BADID123456", "Video 1", youtubeapi.RSSLink{"http://link"}, "", "", mediaGroup}

		videos := []Video{}

		if isEntryInVideoList(&entry, &videos) {
			t.Error("youtubeapi.RSSVideoEntry should not have been found in video list")
		}
	})
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

func TestGetLocalVideosFromDirList(t *testing.T) {
	t.Run("With valid Dirlist", func(t *testing.T) {
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
			MockFileInfo{
				"My Video-basAIdKsyIA.mkv",
				33000000,
				false,
			},
		}

		path := "/test/path"

		expectedVideos := []Video{
			Video{
				path + "/" + dirlist[0].Name(),
				"dCqJ6iPHus0",
				"mp4",
				path,
			},
			Video{
				path + "/" + dirlist[2].Name(),
				"-AC4HwzAK7A",
				"mp4",
				path,
			},
			Video{
				path + "/" + dirlist[4].Name(),
				"_gPsIiKtybA",
				"mp4",
				path,
			},
			Video{
				path + "/" + dirlist[5].Name(),
				"basAIdKsyIA",
				"mkv",
				path,
			},
		}

		videos, err := getLocalVideosFromDirList(&dirlist, path)
		if err != nil {
			t.Errorf("getLocalVideosFromDirList returned an error %s", err)
		}

		if !reflect.DeepEqual(expectedVideos, *videos) {
			t.Errorf("Video list does not match response from function: Expected\n%+v\n, got\n%+v", expectedVideos, *videos)
		}
	})

	t.Run("With invalid Dirlist", func(t *testing.T) {
		dirlist := []os.FileInfo{}
		path := "/test/path"
		videos, err := getLocalVideosFromDirList(&dirlist, path)
		if err != nil {
			t.Errorf("getLocalVideosFromDirList returned an error %s", err)
		}

		if len(*videos) > 0 {
			t.Errorf("getLocalVideosFromDirList with no dir list returned something %s", videos)
		}
	})
}

func TestGetFileType(t *testing.T) {
	t.Run("Standard mp4 extension should come back correct", func(t *testing.T) {
		file := "test.mp4"
		fileType, err := getFileType(file)
		if err != nil {
			t.Errorf("getFileType returned an error %s", err)
		}
		if fileType != "mp4" {
			t.Errorf("getFileType did not return correct result. Expecte mp4, got %s", fileType)
		}
	})

	t.Run("Capitalised mp4 extension should come back in lower case", func(t *testing.T) {
		file := "test.MP4"
		fileType, err := getFileType(file)
		if err != nil {
			t.Errorf("getFileType returned an error %s", err)
		}
		if fileType != "mp4" {
			t.Errorf("getFileType did not return correct result. Expecte mp4, got %s", fileType)
		}
	})

	t.Run("Throws error when dot on end but no file type exists", func(t *testing.T) {
		file := "test."
		_, err := getFileType(file)
		if err == nil {
			t.Errorf("getFileType did not return an error when provided filename %s", file)
		}
	})

	t.Run("Throws error with no extension", func(t *testing.T) {
		file := "test"
		_, err := getFileType(file)
		if err == nil {
			t.Errorf("getFileType did not return an error when provided filename %s", file)
		}
	})

	t.Run("Throws error with no filename", func(t *testing.T) {
		file := ""
		_, err := getFileType(file)
		if err == nil {
			t.Errorf("getFileType did not return an error when provided filename %s", file)
		}
	})
}

func TestIsFileMP4(t *testing.T) {
	t.Run("A file string with an mp4 extension returns true", func(t *testing.T) {
		file := "test.mp4"
		isValid, err := isMP4(file)
		if err != nil {
			t.Errorf("isMP4 returned an error %s", err)
		}
		if !isValid {
			t.Errorf("isMP4 returned false for %s", file)
		}
	})

	t.Run("A file string with a doc extension returns false", func(t *testing.T) {
		file := "test.doc"
		isValid, err := isMP4(file)
		if err != nil {
			t.Errorf("isMP4 returned an error %s", err)
		}
		if isValid {
			t.Errorf("isMP4 returned true for %s", file)
		}
	})

	t.Run("A filename with a dot at the end but no filetype extension returns an error", func(t *testing.T) {
		file := "test."
		_, err := isMP4(file)
		if err == nil {
			t.Error("isMP4 did not return an error")
		}
	})

	t.Run("A filename with no extension returns an error", func(t *testing.T) {
		file := "test"
		_, err := isMP4(file)
		if err == nil {
			t.Error("isMP4 did not return an error")
		}
	})

	t.Run("No filename returns an error", func(t *testing.T) {
		file := ""
		_, err := isMP4(file)
		if err == nil {
			t.Error("isMP4 did not return an error")
		}
	})
}

func TestIsFileMKV(t *testing.T) {
	t.Run("A file string with an mpkv extension returns true", func(t *testing.T) {
		file := "test.mkv"
		isValid, err := isMKV(file)
		if err != nil {
			t.Errorf("isMKV returned an error %s", err)
		}
		if !isValid {
			t.Errorf("isMKV returned false for %s", file)
		}
	})

	t.Run("A file string with a doc extension returns false", func(t *testing.T) {
		file := "test.doc"
		isValid, err := isMKV(file)
		if err != nil {
			t.Errorf("isMKV returned an error %s", err)
		}
		if isValid {
			t.Errorf("isMKV returned true for %s", file)
		}
	})

	t.Run("A filename with a dot at the end but no filetype extension returns an error", func(t *testing.T) {
		file := "test."
		_, err := isMKV(file)
		if err == nil {
			t.Error("isMKV did not return an error")
		}
	})

	t.Run("A filename with no extension returns an error", func(t *testing.T) {
		file := "test"
		_, err := isMKV(file)
		if err == nil {
			t.Error("isMKV did not return an error")
		}
	})

	t.Run("No filename returns an error", func(t *testing.T) {
		file := ""
		_, err := isMKV(file)
		if err == nil {
			t.Error("isMKV did not return an error")
		}
	})
}
