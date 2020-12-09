package main

import (
	"fmt"
	"hyperfocus.systems/youtube-curator-server/collection"
	"hyperfocus.systems/youtube-curator-server/videometadata"
	"hyperfocus.systems/youtube-curator-server/videometadata/tageditor"
)

func main() {
	ytc := collection.Feeds["Ashens"]
	videos, err := collection.GetLocalVideosByYTChannel(&ytc)
	if err != nil {
		panic(err)
	}

	path := (*videos)[0].Path
	fmt.Println(path)

	metadataProvider := &tageditor.MP4MetadataCommandProvider{}
	_, err = videometadata.GetVideoMetadata(path, metadataProvider)

	if err != nil {
		panic(err)
	}

	fmt.Println(metadataProvider)
}
