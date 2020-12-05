package main

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/collection"
	// "hyperfocus.systems/youtube-curator-server/videometadata"
	"hyperfocus.systems/youtube-curator-server/videometadata/mkvinfo"
)

func main() {
	ytc := collection.Feeds["NickRobinson"]
	videos, err := collection.GetLocalVideosByYTChannel(&ytc)
	if err != nil {
		panic(err)
	}

	path := (*videos)[0].Path
	fmt.Println(path)

	metadataProvider := &mkvinfo.MKVMetadataProvider{}
	_, err = metadataProvider.LoadVideoFile(path)

	if err != nil {
		panic(err)
	}

	fmt.Println(metadataProvider)
}
