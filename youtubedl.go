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
	return fmt.Sprintf("youtube-dl -f best -i --write-description --add-metadata --all-subs --sub-format \"srt\" --embed-subs -o \"/media/Drive/Videos/Youtube/%s/%%(title)s-%%(id)s.%%(ext)s\" %s", chann.Name, str)
}

func getYoutubeDLCommandForVideoList(chann *YTChannel, list *[]VideoEntry) string {
	var youtubeDlList []string
	for _, entry := range *list {
		youtubeDlList = append(youtubeDlList, "\""+entry.Link.Href+"\"")
	}

	downloadString := strings.Join(youtubeDlList, " ")

	return getYoutubeDLCommandForYTChannel(chann, downloadString)
}

func getCommandForArchivalType(ytchan *YTChannel, videos *[]VideoEntry) (string, error) {
	if ytchan.ArchivalMode == ArchivalModeCurated {
		return getYoutubeDLCommandForVideoList(ytchan, videos), nil
	} else if ytchan.ArchivalMode == ArchivalModeArchive {
		return getYoutubeDLCommandForYTChannel(ytchan, ytchan.ChannelURL), nil
	}

	return "", fmt.Errorf("Archival Type for provided channel is invalid. Got %s from channel %s", ytchan.ArchivalMode, ytchan)
}
