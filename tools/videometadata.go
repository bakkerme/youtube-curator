package main

// VideoMetadataProvider is an interface for providing metadata from a video file
type VideoMetadataProvider interface {
	LoadVideoFile(path string)
	Title() string
	Description() string
	Creator() string
	PublishedAt() string
	Duration() int
	Description() string
}
