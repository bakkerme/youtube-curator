module hyperfocus.systems/youtube-curator-server/tools

go 1.15

require hyperfocus.systems/youtube-curator-server/config v1.0.0

replace hyperfocus.systems/youtube-curator-server/config => ../config

require hyperfocus.systems/youtube-curator-server/youtubeapi v1.0.0

replace hyperfocus.systems/youtube-curator-server/youtubeapi => ../youtubeapi

require hyperfocus.systems/youtube-curator-server/youtubedl v1.0.0

replace hyperfocus.systems/youtube-curator-server/youtubedl => ../youtubedl

require hyperfocus.systems/youtube-curator-server/utils v1.0.0

replace hyperfocus.systems/youtube-curator-server/utils => ../utils

require hyperfocus.systems/youtube-curator-server/collection v1.0.0

replace hyperfocus.systems/youtube-curator-server/collection => ../collection

require hyperfocus.systems/youtube-curator-server/videometadata v1.0.0

replace hyperfocus.systems/youtube-curator-server/videometadata => ../videometadata

require hyperfocus.systems/youtube-curator-server/videometadata/mkvinfo v1.0.0

replace hyperfocus.systems/youtube-curator-server/videometadata/mkvinfo => ../videometadata/mkvinfo

require hyperfocus.systems/youtube-curator-server/videometadata/tageditor v1.0.0

replace hyperfocus.systems/youtube-curator-server/videometadata/tageditor => ../videometadata/tageditor
