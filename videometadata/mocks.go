package videometadata

import (
	"errors"
)

// MockVideoMetadata mocks the Provider system
type MockVideoMetadata struct {
	GetReturn      *Response
	GetReturnError bool
	SetReturnError bool
}

// Get loads up a video file and returns metadata about the video,
// given the correct VideoMetadataCommandProvider for the type of video
func (l MockVideoMetadata) Get(path string, pr CommandProvider) (*Response, error) {
	if l.GetReturnError {
		return nil, errors.New("Something bad happened")
	}

	return l.GetReturn, nil
}

// Set will set metadata on a video file given the correct
// VideoMetadataCommandProvider for the type of video file
func (l MockVideoMetadata) Set(path string, metadata *Metadata, pr CommandProvider) error {
	if l.SetReturnError {
		return errors.New("Something bad happened")
	}

	return nil
}
