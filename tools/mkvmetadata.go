package main

import (
	"fmt"
	"time"
)

// MKVMetadataProvider provides parsing for MKV metadata
type MKVMetadataProvider struct {
	title       string
	description string
	creator     string
	publishedAt string
	duration    time.Duration
}

// LoadVideoFile loads an MKV video file's metadata from disc
func (m *MKVMetadataProvider) LoadVideoFile(path string) (VideoMetadataError, error) {
	resp, err := loadVideoFile(path)
	if err != nil {
		return nil, fmt.Errorf("Filed to parse video file %s, error: %s", path, err)
	}

	m.title = resp.metadata.title
	m.description = resp.metadata.description
	m.creator = resp.metadata.creator
	m.publishedAt = resp.metadata.publishedAt
	m.duration = resp.metadata.duration

	if resp.parseError != nil {
		return resp.parseError, nil
	}

	return nil, nil
}

// Title returns the video title
func (m *MKVMetadataProvider) Title() string {
	return m.title
}

// Description returns the video description
func (m *MKVMetadataProvider) Description() string {
	return m.description
}

// Creator returns the videos creator
func (m *MKVMetadataProvider) Creator() string {
	return m.creator
}

// PublishedAt returns the date string of the published video
func (m *MKVMetadataProvider) PublishedAt() string {
	return m.publishedAt
}

// Duration returns the duration of the video as a time string
func (m *MKVMetadataProvider) Duration() time.Duration {
	return m.duration
}
