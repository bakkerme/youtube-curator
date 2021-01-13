package youtubeapi

// SearchResponseJSON is an example response from the Youtube Search API for testing purposes
var SearchResponseJSON = `
{
  "kind": "youtube#searchListResponse",
  "etag": "JUwDkhuk2MfT9JohAdWLTXV45aM",
  "nextPageToken": "CAUQAA",
  "regionCode": "AU",
  "pageInfo": {
    "totalResults": 522,
    "resultsPerPage": 3
  },
  "items": [
    {
      "kind": "youtube#searchResult",
      "etag": "mXlJ6j7wGlbm9SjkEyvihSwo0tU",
      "id": {
        "kind": "youtube#video",
        "videoId": "18-elPdai_1"
      },
      "snippet": {
        "publishedAt": "2012-10-01T15:27:35Z",
        "channelId": "UCS-WzPVpAAli-1IfEG2lN8A",
        "title": "Test Video New",
        "description": "Test Description New",
        "thumbnails": {
          "default": {
            "url": "https://i.ytimg.com/vi/18-elPdai_1/default.jpg",
            "width": 120,
            "height": 90
          },
          "medium": {
            "url": "https://i.ytimg.com/vi/18-elPdai_1/mqdefault.jpg",
            "width": 320,
            "height": 180
          },
          "high": {
            "url": "https://i.ytimg.com/vi/18-elPdai_1/hqdefault.jpg",
            "width": 480,
            "height": 360
          }
        },
        "channelTitle": "Test Guy",
        "liveBroadcastContent": "none",
        "publishTime": "2012-10-01T15:27:35Z"
      }
    },
    {
      "kind": "youtube#searchResult",
      "etag": "-fusrPK0jUxsR3-7UT7as7j4sGM",
      "id": {
        "kind": "youtube#video",
        "videoId": "OGK8gnP4TfA"
      },
      "snippet": {
        "publishedAt": "2018-12-03T23:20:21Z",
        "channelId": "UCS-WzPVpAAli-1IfEG2lN8A",
        "title": "Test Video 1",
        "description": "Test Description",
        "thumbnails": {
          "default": {
            "url": "https://i.ytimg.com/vi/OGK8gnP4TfA/default.jpg",
            "width": 120,
            "height": 90
          },
          "medium": {
            "url": "https://i.ytimg.com/vi/OGK8gnP4TfA/mqdefault.jpg",
            "width": 320,
            "height": 180
          },
          "high": {
            "url": "https://i.ytimg.com/vi/OGK8gnP4TfA/hqdefault.jpg",
            "width": 480,
            "height": 360
          }
        },
        "channelTitle": "Test Guy",
        "liveBroadcastContent": "none",
        "publishTime": "2018-12-03T23:20:21Z"
      }
    },
    {
      "kind": "youtube#searchResult",
      "etag": "nWfU-BRD9p-BGwQ_oFpSv7YmaeI",
      "id": {
        "kind": "youtube#video",
        "videoId": "FazJqPQ6xSs"
      },
      "snippet": {
        "publishedAt": "2019-06-03T19:00:06Z",
        "channelId": "UCS-WzPVpAAli-1IfEG2lN8A",
        "title": "Test Video 2",
        "description": "Test Description 2",
        "thumbnails": {
          "default": {
            "url": "https://i.ytimg.com/vi/FazJqPQ6xSs/default.jpg",
            "width": 120,
            "height": 90
          },
          "medium": {
            "url": "https://i.ytimg.com/vi/FazJqPQ6xSs/mqdefault.jpg",
            "width": 320,
            "height": 180
          },
          "high": {
            "url": "https://i.ytimg.com/vi/FazJqPQ6xSs/hqdefault.jpg",
            "width": 480,
            "height": 360
          }
        },
        "channelTitle": "Test Guy",
        "liveBroadcastContent": "none",
        "publishTime": "2019-06-03T19:00:06Z"
      }
    }
  ]
}
`

// SearchResponseEmptyJSON is an response from the YouTube API, containing a no videos
var SearchResponseEmptyJSON = `
{
  "kind": "youtube#searchListResponse",
  "etag": "JUwDkhuk2MfT9JohAdWLTXV45aM",
  "items": [],
  "pageInfo": {
    "totalResults": 0,
    "resultsPerPage": 0
  }
}`

// VideoResponseJSON is an response from the YouTube Video API, containing a single video
var VideoResponseJSON = `
{
  "kind": "youtube#videoListResponse",
  "etag": "1-lTCZCHtgPgr709KQ0ef2Mu4oM",
  "items": [
    {
      "kind": "youtube#video",
      "etag": "jB-DuI2TOpg-o1d5hnzty8kExw8",
      "id": "18-elPdai_1",
      "snippet": {
        "publishedAt": "2012-10-01T15:27:35Z",
        "channelId": "UCS-WzPVpAAli-1IfEG2lN8A",
        "title": "Test Video New",
        "description": "Test Description New",
        "thumbnails": {
          "default": {
            "url": "https://i2.ytimg.com/vi/KQA9Na4aOa1/hqdefault.jpg",
            "width": 120,
            "height": 90
          },
          "medium": {
            "url": "https://i.ytimg.com/vi/KQA9Na4aOa1/mqdefault.jpg",
            "width": 320,
            "height": 180
          },
          "high": {
            "url": "https://i.ytimg.com/vi/KQA9Na4aOa1/hqdefault.jpg",
            "width": 480,
            "height": 360
          },
          "standard": {
            "url": "https://i.ytimg.com/vi/KQA9Na4aOa1/sddefault.jpg",
            "width": 640,
            "height": 480
          },
          "maxres": {
            "url": "https://i.ytimg.com/vi/KQA9Na4aOa1/maxresdefault.jpg",
            "width": 1280,
            "height": 720
          }
        },
        "channelTitle": "Test Guy",
        "tags": [
          "IMA",
          "TEST",
          "TAG"
        ],
        "categoryId": "22",
        "liveBroadcastContent": "none",
        "defaultLanguage": "en",
        "localized": {
          "title": "Test Video New",
          "description": "Test Description New"
        },
        "defaultAudioLanguage": "en"
      }
    }
  ],
  "pageInfo": {
    "totalResults": 1,
    "resultsPerPage": 1
  }
}
`

// PlaylistItemsResponseJSON is an example response from the playlist Youtube API
var PlaylistItemsResponseJSON = `
{
  "kind": "youtube#playlistItemListResponse",
  "etag": "c6a05HjuPsmPxbDxMmt-196SvPI",
  "nextPageToken": "CAUQAA",
  "items": [
    {
      "kind": "youtube#playlistItem",
      "etag": "b_mpir6DuPAQXBP2QtAT9y2Dp8c",
      "id": "UExOejRVbjkycEdOeFE5dk5nbW5DeDdkd2NoUEpHSjNJUS41NTZEOThBNThFOUVGQkVB",
      "snippet": {
        "publishedAt": "2020-09-30T19:22:21Z",
        "channelId": "UCjU-Cwjfqbo2hMRItlXwnnQ",
        "title": "Test Video New",
        "description": "Test Description New",
        "thumbnails": {
          "default": {
            "url": "https://i.ytimg.com/vi/GhMOw141DKg/default.jpg",
            "width": 120,
            "height": 90
          },
          "medium": {
            "url": "https://i.ytimg.com/vi/GhMOw141DKg/mqdefault.jpg",
            "width": 320,
            "height": 180
          },
          "high": {
            "url": "https://i.ytimg.com/vi/GhMOw141DKg/hqdefault.jpg",
            "width": 480,
            "height": 360
          },
          "standard": {
            "url": "https://i.ytimg.com/vi/GhMOw141DKg/sddefault.jpg",
            "width": 640,
            "height": 480
          }
        },
        "channelTitle": "Test Guy",
        "playlistId": "PLNz4Un92pGNxQ9vNgmnCx7dwchPJGJ3IQ",
        "position": 0,
        "resourceId": {
          "kind": "youtube#video",
          "videoId": "GhMOw141DKg"
        }
      }
    },
    {
      "kind": "youtube#playlistItem",
      "etag": "Nzlqatt6khGGAufkcqgLhEZ6fbg",
      "id": "UExOejRVbjkycEdOeFE5dk5nbW5DeDdkd2NoUEpHSjNJUS42Qzk5MkEzQjVFQjYwRDA4",
      "snippet": {
        "publishedAt": "2020-03-30T00:22:24Z",
        "channelId": "UCjU-Cwjfqbo2hMRItlXwnnQ",
        "title": "Test Video 1",
        "description": "Test Description 1",
        "thumbnails": {
          "default": {
            "url": "https://i.ytimg.com/vi/YEa2aj9KYQA/default.jpg",
            "width": 120,
            "height": 90
          },
          "medium": {
            "url": "https://i.ytimg.com/vi/YEa2aj9KYQA/mqdefault.jpg",
            "width": 320,
            "height": 180
          },
          "high": {
            "url": "https://i.ytimg.com/vi/YEa2aj9KYQA/hqdefault.jpg",
            "width": 480,
            "height": 360
          },
          "standard": {
            "url": "https://i.ytimg.com/vi/YEa2aj9KYQA/sddefault.jpg",
            "width": 640,
            "height": 480
          }
        },
        "channelTitle": "Test Guy",
        "playlistId": "PLNz4Un92pGNxQ9vNgmnCx7dwchPJGJ3IQ",
        "position": 1,
        "resourceId": {
          "kind": "youtube#video",
          "videoId": "YEa2aj9KYQA"
        }
      }
    },
    {
      "kind": "youtube#playlistItem",
      "etag": "5LQBWt1kwMx7Uqtb2gTHP0bnHOw",
      "id": "UExOejRVbjkycEdOeFE5dk5nbW5DeDdkd2NoUEpHSjNJUS4zMEQ1MEIyRTFGNzhDQzFB",
      "snippet": {
        "publishedAt": "2020-01-23T16:36:27Z",
        "channelId": "UCjU-Cwjfqbo2hMRItlXwnnQ",
        "title": "Test Video 2",
        "description": "Test Description 2",
        "thumbnails": {
          "default": {
            "url": "https://i.ytimg.com/vi/fZYBhmteJDE/default.jpg",
            "width": 120,
            "height": 90
          },
          "medium": {
            "url": "https://i.ytimg.com/vi/fZYBhmteJDE/mqdefault.jpg",
            "width": 320,
            "height": 180
          },
          "high": {
            "url": "https://i.ytimg.com/vi/fZYBhmteJDE/hqdefault.jpg",
            "width": 480,
            "height": 360
          },
          "standard": {
            "url": "https://i.ytimg.com/vi/fZYBhmteJDE/sddefault.jpg",
            "width": 640,
            "height": 480
          }
        },
        "channelTitle": "Test Guy",
        "playlistId": "PLNz4Un92pGNxQ9vNgmnCx7dwchPJGJ3IQ",
        "position": 2,
        "resourceId": {
          "kind": "youtube#video",
          "videoId": "fZYBhmteJDE"
        }
      }
    }
  ],
  "pageInfo": {
    "totalResults": 54,
    "resultsPerPage": 3
  }
}
`

// VideoResponseRSSXML is an example repsonse from the Youtube RSS update system
var VideoResponseRSSXML = `
<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns:yt="http://www.youtube.com/xml/schemas/2015" xmlns:media="http://search.yahoo.com/mrss/" xmlns="http://www.w3.org/2005/Atom">
	<link rel="self" href="http://www.youtube.com/feeds/videos.xml?channel_id=UCS-WzPVpAAli-1IfEG2lN8A"/>
	<id>yt:channel:UCS-WzPVpAAli-1IfEG2lN8A</id>
	<yt:channelId>UCS-WzPVpAAli-1IfEG2lN8A</yt:channelId>
	<title>Test Guy</title>
	<link rel="alternate" href="https://www.youtube.com/channel/UCS-WzPVpAAli-1IfEG2lN8A"/>
	<author>
		<name>Test Guy</name>
		<uri>https://www.youtube.com/channel/UCS-WzPVpAAli-1IfEG2lN8A</uri>
	</author>
	<published>2010-01-30T19:58:04+00:00</published>
	<entry>
		<id>yt:video:KQA9Na4aOa1</id>
		<yt:videoId>KQA9Na4aOa1</yt:videoId>
		<yt:channelId>UCS-WzPVpAAli-1IfEG2lN8A</yt:channelId>
		<title>Test Video New</title>
		<link rel="alternate" href="https://www.youtube.com/watch?v=KQA9Na4aOa1"/>
		<author>
			<name>Test Guy</name>
			<uri>https://www.youtube.com/channel/UCS-WzPVpAAli-1IfEG2lN8A</uri>
		</author>
		<published>2020-11-06T19:00:01+00:00</published>
		<updated>2020-11-06T23:12:15+00:00</updated>
		<media:group>
			<media:title>Test Video 1</media:title>
			<media:content url="https://www.youtube.com/v/KQA9Na4aOa1?version=3" type="application/x-shockwave-flash" width="640" height="390"/>
			<media:thumbnail url="https://i2.ytimg.com/vi/KQA9Na4aOa1/hqdefault.jpg" width="480" height="360"/>
			<media:description>Test Description New</media:description>
			<media:community>
				<media:starRating count="470" average="4.97" min="1" max="5"/>
				<media:statistics views="4602"/>
			</media:community>
		</media:group>
	</entry>
	<entry>
		<id>yt:video:OGK8gnP4TfA</id>
		<yt:videoId>OGK8gnP4TfA</yt:videoId>
		<yt:channelId>UCS-WzPVpAAli-1IfEG2lN8A</yt:channelId>
		<title>Test Video 1</title>
		<link rel="alternate" href="https://www.youtube.com/watch?v=OGK8gnP4TfA"/>
		<author>
			<name>Test Guy</name>
			<uri>https://www.youtube.com/channel/UCS-WzPVpAAli-1IfEG2lN8A</uri>
		</author>
		<published>2020-11-06T19:00:01+00:00</published>
		<updated>2020-11-06T23:12:15+00:00</updated>
		<media:group>
			<media:title>Test Video 1</media:title>
			<media:content url="https://www.youtube.com/v/OGK8gnP4TfA?version=3" type="application/x-shockwave-flash" width="640" height="390"/>
			<media:thumbnail url="https://i2.ytimg.com/vi/OGK8gnP4TfA/hqdefault.jpg" width="480" height="360"/>
			<media:description>Test Description</media:description>
			<media:community>
				<media:starRating count="470" average="4.97" min="1" max="5"/>
				<media:statistics views="4602"/>
			</media:community>
		</media:group>
	</entry>
	<entry>
		<id>yt:video:FazJqPQ6xSs</id>
		<yt:videoId>FazJqPQ6xSs</yt:videoId>
		<yt:channelId>UCS-WzPVpAAli-1IfEG2lN8A</yt:channelId>
		<title>Test Video 2</title>
		<link rel="alternate" href="https://www.youtube.com/watch?v=FazJqPQ6xSs"/>
		<author>
			<name>Test Guy</name>
			<uri>https://www.youtube.com/channel/UCS-WzPVpAAli-1IfEG2lN8A</uri>
		</author>
		<published>2020-11-06T19:00:01+00:00</published>
		<updated>2020-11-06T23:12:15+00:00</updated>
		<media:group>
			<media:title>Test Video 2</media:title>
			<media:content url="https://www.youtube.com/v/FazJqPQ6xSs?version=3" type="application/x-shockwave-flash" width="640" height="390"/>
			<media:thumbnail url="https://i2.ytimg.com/vi/FazJqPQ6xSs/hqdefault.jpg" width="480" height="360"/>
			<media:description>Test Description 2</media:description>
			<media:community>
				<media:starRating count="470" average="4.97" min="1" max="5"/>
				<media:statistics views="4602"/>
			</media:community>
		</media:group>
	</entry>
</feed>
`
