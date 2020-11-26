package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// MKVMetadataProvider provides parsing for MKV metadata
type MKVMetadataProvider struct {
	Title       string
	Description string
	Creator     string
	PublishedAt string
	Duration    time.Duration
}

// LoadVideoFile loads an MKV video file's metadata from disc
func (m *MKVMetadataProvider) LoadVideoFile(path string) error {
	cmd := exec.Command(fmt.Sprintf("mkvinfo \"%s\"", path))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Could not load metadata for %s. Error %s", path, err)
	}
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Could not load metadata for %s. Failed on output stage. Error %s", path, err)
	}

	err := m.parseMKVInfoOutput(string(out))
	if err != nil {
		return fmt.Errorf("Could not load metadata for %s. Failed on parse stage. Error %s", path, err)
	}

	return nil
}

// Title returns the video title
func (m *MKVMetadataProvider) Title() string {
	return m.Title
}

// Description returns the video description
func (m *MKVMetadataProvider) Description() string {
	return m.Description
}

// Creator returns the videos creator
func (m *MKVMetadataProvider) Creator() string {
	return m.Creator
}

// PublishedAt returns the date string of the published video
func (m *MKVMetadataProvider) PublishedAt() string {
	return m.PublishedAt
}

// Duration returns the duration of the video as a time string
func (m *MKVMetadataProvider) Duration() string {
	return m.Duration
}

func (m *MKVMetadataProvider) parseMKVInfoOutput(output string) error {
	m.Title = parseDescription(output)
	m.Description = parseDescription(output)
	m.Creator = parseCreator(output)
	m.PublishedAt = parsePublishedAt(output)

	duration, err = parseDuration(output)
	if err != nil {
		return err
	}
	m.Duration = duration

	return nil
}

func parseTitle(output string) string {
	return parseOutputStringForRegex(`(?m)\| \+ Title: (.*$)`, output)
}

func parseDescription(output string) string {
	return parseOutputStringForRegex(`(?msU)\+ Name: DESCRIPTION.*String: (.*)^\|`, output)
}

func parseCreator(output string) string {
	return parseOutputStringForRegex(`(?msU)ARTIST.*String: (.*^)`, output)
}

func parsePublishedAt(output string) string {
	return parseOutputStringForRegex(`(?msU)DATE.*String: (.*^)`, output)
}

func parseDuration(output string) (time.Duration, error) {
	str := parseOutputStringForRegex(`(?msU)Duration: (.*^)`, output)
	milisecondless := strings.Split(str, ".")[0]
	units := strings.Split(milisecondless, ':')
	durationString := fmt.Sprintf("%sh%sm%ss")
	duration, err := time.ParseDuration(duration)

	if err != nil {
		return nil, fmt.Errorf("Could not parse MKV duration. Attempted to parse %s. Error %s", str, err)
	}

	return duration, nil
}

func parseOutputStringForRegex(regex string, output string) {
	re := regexp.MustCompile(regex)
	return re.FindAllString(output, -1)[0]
}
