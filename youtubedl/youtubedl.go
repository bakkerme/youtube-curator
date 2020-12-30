package youtubedl

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
	"strings"
)

var youtubeDLCommand = []string{
	"youtube-dl",
	"--format",
	"\"(bestvideo[vcodec^=avc1][height=1080][fps>30]/bestvideo[vcodec^=avc1][height=1080]/bestvideo[vcodec^=avc1][height=720][fps>30]/bestvideo[vcodec^=avc1][height=720]/bestvideo[vcodec^=avc1][height=480][fps>30]/bestvideo[vcodec^=avc1][height=480]/bestvideo[vcodec^=avc1][height=360][fps>30]/bestvideo[vcodec^=avc1][height=360]/bestvideo[vcodec^=avc1][height=240][fps>30]/bestvideo[vcodec^=avc1][height=240]/bestvideo[vcodec^=avc1][height=144][fps>30]/bestvideo[vcodec^=avc1][height=144]/bestvideo[vcodec^=avc1])+(bestaudio[acodec^=mp4a]/bestaudio)/best\"",
	"--verbose",
	"--force-ipv4",
	"--sleep-interval 5",
	"--max-sleep-interval 30",
	"--ignore-errors",
	"--no-continue",
	"--no-overwrites",
	"--download-archive archive.log",
	"--add-metadata",
	"--all-subs",
	"--sub-format \"srt\"",
	"--embed-subs",
	"--output \"%(upload_date)s - %(title)s-%(id)s.%(ext)s\"",
	"--merge-output-format \"mkv\"",
}

func getYoutubeDLCommandForYTChannel(ytchan collection.YTChannel, str string, baseDir string) string {
	cdCommand := fmt.Sprintf("cd %s%s;", baseDir, ytchan.Name())
	return fmt.Sprintf("%s %s %s", cdCommand, strings.Join(youtubeDLCommand, " "), str)
}

func getYoutubeDLCommandForVideoList(ytchan collection.YTChannel, list *[]youtubeapi.RSSVideoEntry, baseDir string) string {
	var youtubeDlList []string
	for _, entry := range *list {
		youtubeDlList = append(youtubeDlList, "\""+entry.Link.Href+"\"")
	}

	downloadString := strings.Join(youtubeDlList, " ")

	return getYoutubeDLCommandForYTChannel(ytchan, downloadString, baseDir)
}

// GetCommandForArchivalType provides a YoutubeDL command for a YTChannel to download a number of VideoEntrys
func GetCommandForArchivalType(ytchan collection.YTChannel, videos *[]youtubeapi.RSSVideoEntry, cf *config.Config) (string, error) {
	if ytchan.ArchivalMode() == collection.ArchivalModeCurated {
		return getYoutubeDLCommandForVideoList(ytchan, videos, cf.VideoDirPath), nil
	} else if ytchan.ArchivalMode() == collection.ArchivalModeArchive {
		return getYoutubeDLCommandForYTChannel(ytchan, ytchan.ChannelURL(), cf.VideoDirPath), nil
	}

	return "", fmt.Errorf("Archival Type for provided channel is invalid. Got %s from channel %s", ytchan.ArchivalMode(), ytchan)
}
