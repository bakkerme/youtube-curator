package main

import (
	"fmt"
	"strings"
)

var youtubeDLCommand = []string{
	"youtube-dl",
	"-f best",
	"-i",
	"--write-description",
	"--add-metadata",
	"--all-subs",
	"--sub-format \"srt\"",
	"--embed-subs",
	"-o \"/media/Drive/Videos/Youtube/%s/%%(title)s-%%(id)s.%%(ext)s %s\"",
}

func getYoutubeDLCommandForYTChannel(chann *YTChannel, str string) string {
	return fmt.Sprintf("youtube-dl -f best -i --write-description --add-metadata --all-subs --sub-format \"srt\" --embed-subs -o \"/media/Drive/Videos/Youtube/%s/%%(title)s-%%(id)s.%%(ext)s %s\"", chann.Name, str)
}

func getYoutubeDLCommandForVideoList(chann *YTChannel, list *[]VideoEntry) string {
	var youtubeDlList []string
	for _, entry := range *list {
		youtubeDlList = append(youtubeDlList, "\""+entry.Link.Href+"\"")
	}

	downloadString := strings.Join(youtubeDlList, ", ")

	return getYoutubeDLCommandForYTChannel(chann, downloadString)
}
