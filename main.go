package main

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
	"errors"
	"log"
)

func main() {
	file, err := ioutil.ReadFile("./tests/MichaelMJD.xml")
	if err != nil {
		log.Fatalf("Loading RSS feed xml failed: %s", err)
	}

	var rss RSS
	if err := xml.Unmarshal([]byte(file), &rss); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}

	log.Println(rss.ID)
	log.Println(rss.Title)
	log.Println(rss.Entry)
}

func getVideoIDFromFileName(filename string) (string, error) {
	parseError := errors.New("Could not parse video ID")

	splits := strings.Split(filename, ".")
	withoutType := splits[0]
	id := withoutType[len(withoutType)-11:len(withoutType)] // Get last 11 chars

	if id == "" {
		return "", parseError
	}

	if len(id) != len(strings.ReplaceAll(id, " ", "")) { // This is probably not an ID
		return "", parseError
	}

	return id, nil
}

func isEntryInVideoList(entry *Entry, videos *[]Video) bool {
	match := false
	for _, video := range *videos {
		if video.ID == entry.ID {
			match = true
		}
	}

	return match;
}

func getEntriesNotInVideoList(entries *[]Entry, videos *[]Video) *[]Entry {
	var notInVideoList []Entry
	for _, entry := range *entries {
		match := isEntryInVideoList(&entry, videos)
		if !match { // Entry isn't in our list of videos
			notInVideoList = append(notInVideoList, entry)
		}
	}

	return &notInVideoList
}

func getLocalVideosByChannel(channel *Channel) *[]Video {
	path := "/media/Drive/Videos/Youtube/" + channel.name
	dirlist, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("Could not get videos off disk for channel %s, error %s", channel.name, err)
	}

	var videos []Video
	for _, file := range files {
		videoPath = path + "/" + file.Name()
		fmt.Println(file.Name())
		video := Video{
			videoPath,
			getVideoIDFromFileName(videoPath),
		}

		append(videos, video)
	}

	return videos
}
