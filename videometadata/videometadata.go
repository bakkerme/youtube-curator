package videometadata

import (
	"time"
)

// Provider provides an interface for getting and setting Video metadata
type Provider interface {
	Get(path string, pr CommandProvider) (*Response, error)
	Set(path string, metadata *Metadata, pr CommandProvider) error
}

// CommandProvider is an interface for different providers of metadata. This
// can be used to implement metadata parsers for various types of video container
// formats
type CommandProvider interface {
	Run(string) (string, error)
	ParseTitle(string) (string, error)
	ParseDescription(string) (string, error)
	ParseCreator(string) (string, error)
	ParsePublishedAt(string) (*time.Time, error)
	ParseDuration(string) (*time.Duration, error)
	Set(string, *Metadata) error
}

// Error represents a parse error. This is not necessarily fatal
// and may be expected, depending on the data in the Video's metadata
type Error interface {
	Error() string
	UnparsedFields() []string
}

// Response contains the videos metadata, and parse errors if the video could not be parsed fully
type Response struct {
	Metadata   *Metadata
	ParseError *ParseError
}

// Metadata represents the metadata for the loaded video
type Metadata struct {
	Title       string
	Description string
	Creator     string
	PublishedAt *time.Time
	Duration    *time.Duration
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

// VideoMetadata provides the ability to Get and Set VideoMetadata
type VideoMetadata struct{}

// Get loads up a video file and returns metadata about the video,
// given the correct VideoMetadataCommandProvider for the type of video
func (l VideoMetadata) Get(path string, pr CommandProvider) (*Response, error) {
	out, err := pr.Run(path)
	if err != nil {
		return nil, err
	}

	return buildVideoMetadataResponse(out, path, pr), nil
}

// Set will set metadata on a video file given the correct
// VideoMetadataCommandProvider for the type of video file
func (l VideoMetadata) Set(path string, metadata *Metadata, pr CommandProvider) error {
	return pr.Set(path, metadata)
}
