package mkvinfo

import (
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/videometadata"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// MKVMetadataCommandProvider provides the ability for the VideoMetadata
type MKVMetadataCommandProvider struct{}

// Run MKVInfo on the provided file
func (m MKVMetadataCommandProvider) Run(path string) (string, error) {
	out, err := exec.Command("/usr/bin/mkvinfo", path).Output()
	if err != nil {
		return "", fmt.Errorf("Could not load metadata for %s.\nResponse from mkvinfo was: %s\nError %s", path, out, err)
	}

	return string(out), nil
}

// ParseTitle parses the title from the MKVInfo output
func (m MKVMetadataCommandProvider) ParseTitle(output string) (string, error) {
	return parseOutputStringForRegex(`(?m)\| \+ Title: (?P<content>.*$)`, output)
}

// ParseDescription parses the description from the MKVInfo output
func (m MKVMetadataCommandProvider) ParseDescription(output string) (string, error) {
	return parseOutputStringForRegex(`(?msU)\+ Name: DESCRIPTION.*String: (?P<content>.*)\n^\|`, output)
}

// ParseCreator parses the creator from the MKVInfo output
func (m MKVMetadataCommandProvider) ParseCreator(output string) (string, error) {
	return parseOutputStringForRegex(`(?msU)ARTIST.*String: (?P<content>.*$)`, output)
}

// ParsePublishedAt parses the publishedAt from the MKVInfo output
func (m MKVMetadataCommandProvider) ParsePublishedAt(output string) (*time.Time, error) {
	str, err := parseOutputStringForRegex(`(?msU)DATE.*String: (?P<content>.*$)`, output)
	if err != nil {
		return nil, err
	}

	tm, err := time.Parse("20060102", str)
	if err != nil {
		return nil, err
	}

	return &tm, nil
}

// ParseDuration parses the duration from the MKVInfo output
func (m MKVMetadataCommandProvider) ParseDuration(output string) (*time.Duration, error) {
	str, err := parseOutputStringForRegex(`(?msU)Duration: (?P<content>.*$)`, output)
	if err != nil {
		return nil, err
	}

	milisecondless := strings.Split(str, ".")[0]
	units := strings.Split(milisecondless, ":")
	durationString := fmt.Sprintf("%sh%sm%ss", units[0], units[1], units[2])
	duration, err := time.ParseDuration(durationString)

	if err != nil {
		return nil, err
	}

	return &duration, nil
}

// SetMetadata sets metadata on an mkv item
func (m MKVMetadataCommandProvider) SetMetadata(metadata videometadata.Metadata, path string) error {
	return errors.New("MKVInfo does not implement SetMetadata right now")
}

func parseOutputStringForRegex(regex string, output string) (string, error) {
	re := regexp.MustCompile(regex)
	matches := re.FindStringSubmatch(output)

	if matches == nil {
		return "", fmt.Errorf("Failed to find content for regex %s", regex)
	}

	return matches[1], nil
}
