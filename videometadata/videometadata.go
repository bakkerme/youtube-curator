package videometadata

import (
	"time"
)

// Provider is an interface for providing metadata from a video file
type Provider interface {
	LoadVideoFile(path string)
	Title() string
	Description() string
	Creator() string
	PublishedAt() string
	Duration() int
}

// CommandProvider is an interface for different providers of metadata. This
// can be used to implement metadata parsers for various types of video container
// formats
type CommandProvider interface {
	Run(path string) (string, error)
	ParseTitle(string) (string, error)
	ParseDescription(string) (string, error)
	ParseCreator(string) (string, error)
	ParsePublishedAt(string) (string, error)
	ParseDuration(string) (time.Duration, error)
}

// Error represents a parse error. This is not necessarily fatal
// and may be expected, depending on the data in the Video's metadata
type Error interface {
	Error() string
	UnparsedFields() []string
}

// Response contains the videos metadata, and parse errors if the video could not be parsed fully
type Response struct {
	metadata   *Metadata
	parseError *ParseError
}

// Metadata represents the metadata for the loaded video
type Metadata struct {
	title       string
	description string
	creator     string
	publishedAt string
	duration    time.Duration
}

// ParseError represents an error string and a list of fields that could not be parsed from the video
type ParseError struct {
	err            string
	unparsedFields []string
}

func (pErr *ParseError) Error() string {
	return pErr.err
}

// UnparsedFields returns fields that were unable to be parsed by the parsing
// system
func (pErr *ParseError) UnparsedFields() []string {
	return pErr.unparsedFields
}

// GetVideoMetadata loads up a video file and returns metadata about the video,
// given the correct VideoMetadataCommandProvider for the type of video
func GetVideoMetadata(path string, pr CommandProvider) (*Response, error) {
	out, err := pr.Run(path)
	if err != nil {
		return nil, err
	}

	return buildVideoMetadataResponse(out, path, pr), nil
}
