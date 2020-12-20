package main

import (
	"encoding/json"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/config"
)

func main() {
	cfg, err := config.GetConfig(&config.EnvarConfigProvider{})
	if err != nil {
		panic(err)
	}

	ytChannels, err := collection.GetAvailableYTChannels(cfg)
	if err != nil {
		panic(err)
	}

	channelJSON, err := json.MarshalIndent(ytChannels, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(channelJSON))
}
