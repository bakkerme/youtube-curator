package main

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/config"
	"hyperfocus.systems/youtube-curator-server/youtubeapi"
	"log"
)

func main() {
	config, err := config.GetConfig(&config.EnvarConfigProvider{})
	if err != nil {
		log.Panicf("Config Loader threw an error %s", err)
	}

	id := "UiS27feX8o0"
	video, err := youtubeapi.GetVideoInfo(id, config)
	if err != nil {
		panic(fmt.Sprintf("Could not get video of id %s, %s", id, err))
	}

	fmt.Println(video)
}
