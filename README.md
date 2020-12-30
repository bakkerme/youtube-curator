# Youtube Curator

## Dependencies
* youtube-dl
* ffmpeg
* oapi-codegen

## Setup

Setup .envrc.local with a YOUTUBE_API_KEY and VIDEO_DIR_PATH.

Create folders in the VIDEO_DIR_PATH for each Youtube Channel. Add a config.json with something like the following:

```
{
  "name": "FolderName",
  "rssURL": "https://www.youtube.com/feeds/videos.xml?channel_id=CHANNEL_ID",
  "channelURL": "https://www.youtube.com/channel/CHANNEL_ID",
  "archivalMode": "curated"
}
```

ArchivalMode can be curated or archive.

## TODO
* Refactor out RSS feed tech and replace it with more Youtube API stuff
* Unite all the disparate Video representations
