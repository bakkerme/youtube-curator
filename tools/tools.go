package main

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
)

func main() {
	cfg, err := config.GetConfig(&config.EnvarConfigProvider{})
	if err != nil {
		panic(err)
	}

	ytChannels, err := collection.GetAllLocalVideos(cfg)
	if err != nil {
		panic(err)
	}

	for _, ytc := range *ytChannels {
		fmt.Println(ytc.Path)
	}
}
