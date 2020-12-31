# Youtube Curator

## Dependencies
* youtube-dl
* ffmpeg

## Setup

Setup .envrc.local with a YOUTUBE_API_KEY and VIDEO_DIR_PATH, or add a config file at ~/.config/yt-up2date/config.json:

```
{
  "YoutubeAPIKey": "YOUTUBE_API_KEY goes here, see API dashboard",
  "VideoDirPath": "/your/video/storage/path"
}
```

Create folders in the Vide Dir Path for each Youtube Channel. Add a config.json with something like the following:

```
{
  "name": "FolderName",
  "rssURL": "https://www.youtube.com/feeds/videos.xml?channel_id=CHANNEL_ID",
  "channelURL": "https://www.youtube.com/channel/CHANNEL_ID",
  "archivalMode": "curated"
}
```

ArchivalMode can be "curated" or "archive".

Run:
`go generate`

## TODO
* Initial implementation of API Video lookup functions
* Reimplement Up2Date functionality with Youtube API
* Job function WIP
* Unite all the disparate Video representations
* Refactor disk lookups for more speed
