module hyperfocus.systems/youtube-curator-server/tools

go 1.15

require hyperfocus.systems/youtube-curator-server/config v1.0.0

replace hyperfocus.systems/youtube-curator-server/config => ../config

require hyperfocus.systems/youtube-curator-server/youtubeapi v1.0.0

replace hyperfocus.systems/youtube-curator-server/youtubeapi => ../youtubeapi

require hyperfocus.systems/youtube-curator-server/utils v1.0.0

replace hyperfocus.systems/youtube-curator-server/utils => ../utils
