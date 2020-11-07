package main

// Channel is a struct that represents the URL for each channel
type Channel struct {
	name string
	rssURL string
	channelURL string
	archivalMode string
}

type Video struct {
	Path string
	ID string
}

const archivalModeArchive = "archive"
const archivalModeCurated = "curated"

var feeds = map[string]Channel{
	"65scribe": Channel{
		"65scribe",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC8dJOqcjyiA9Zo9aOxxiCMw",
		"https://www.youtube.com/user/65scribe",
		archivalModeArchive,
	},
	"ashens": Channel{
		"ashens",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCxt9Pvye-9x_AIcb1UtmF1Q",
		"https://www.youtube.com/user/ashens",
		archivalModeArchive,
	},
	"BryanLunduke": Channel{
		"BryanLunduke",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCkK9UDm_ZNrq_rIXCz3xCGA",
		"https://www.youtube.com/user/bryanlunduke",
		archivalModeArchive,
	},
	"DanBell": Channel{
		"DanBell",
		"https://www.youtube.com/feeds/videos.xml?playlist_id=PLNz4Un92pGNxQ9vNgmnCx7dwchPJGJ3IQ",
		"https://www.youtube.com/user/moviedan",
		archivalModeArchive,
	},
	"LinusTechTips": Channel{
		"LinusTechTips",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCXuqSBlHAE6Xw-yeJA0Tunw",
		"https://www.youtube.com/user/LinusTechTips",
		archivalModeCurated,
	},
	"LukeSmith": Channel{
		"LukeSmith",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UC2eYFnH61tmytImy1mTYvhA",
		"https://www.youtube.com/channel/UC2eYFnH61tmytImy1mTYvhA",
		archivalModeArchive,
	},
	"Mario64BetaArchive": Channel{
		"Mario64BetaArchive",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCSar92RCPocysNbvBG84Mxw",
		"https://www.youtube.com/channel/UCSar92RCPocysNbvBG84Mxw",
		archivalModeArchive,
	},
	"Memospore": Channel{
		"Memospore",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UChbm-JCx_jii5xn-2f5nwIA",
		"https://www.youtube.com/user/memospore",
		archivalModeArchive,
	},
	"MichaelMJD": Channel{
		"MichaelMJD",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCS-WzPVpAAli-1IfEG2lN8A",
		"https://www.youtube.com/user/mjd7999",
		archivalModeArchive,
	},
	"RedLetterMedia": Channel{
		"RedLetterMedia",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCrTNhL_yO3tPTdQ5XgmmWjA",
		"https://www.youtube.com/user/RedLetterMedia",
		archivalModeCurated,
	},
	"SurviveTheJive": Channel{
		"SurviveTheJive",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCZAENaOaceQUMd84GDc26EA",
		"https://www.youtube.com/user/ThomasRowsell",
		archivalModeArchive,
	},
	"TechnologyConnections": Channel{
		"TechnologyConnections",
		"https://www.youtube.com/feeds/videos.xml?channel_id=UCy0tKL1T7wFoYcxCe0xjN6Q",
		"https://www.youtube.com/channel/UCy0tKL1T7wFoYcxCe0xjN6Q",
		archivalModeArchive,
	},
}
