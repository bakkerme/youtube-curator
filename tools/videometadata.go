package main

// VideoMetadataProvider is an interface for providing metadata from a video file
type VideoMetadataProvider interface {
	LoadVideoFile(path string)
	Title() string
	Description() string
	Creator() string
	PublishedAt() string
	Duration() int
}

// VideoMetadataError represents a parse error. This is not necessarily fatal
// and may be expected, depending on the data in the Video's metadata
type VideoMetadataError interface {
	Error() string
	UnparsedFields() []string
}
