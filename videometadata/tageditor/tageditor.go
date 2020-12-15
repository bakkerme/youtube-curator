package tageditor

import (
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/utils"
	"hyperfocus.systems/youtube-curator-server/videometadata"
	"regexp"
	"strings"
	"time"
)

// MP4MetadataCommandProvider provides the ability for the VideoMetadata
type MP4MetadataCommandProvider struct{}

// Run MP4Info on the provided file
func (m MP4MetadataCommandProvider) Run(path string) (string, error) {
	return getMetadata(path, &utils.OSCommand{})
}

func getMetadata(path string, osc utils.OSCommandProvider) (string, error) {
	out, err := osc.Run("tageditor", "get", "-f", path)
	if err != nil {
		return "", fmt.Errorf("Could not load metadata for %s.\nResponse from tageditor was: %s\nError %s", path, out, err)
	}

	infoOut, err := osc.Run("tageditor", "info", "-f", path)
	if err != nil {
		return "", fmt.Errorf("Could not load metadata for %s.\nResponse from tageditor was: %s\nError %s", path, out, err)
	}

	return string(*out) + string(*infoOut), nil
}

// ParseTitle parses the title from the MP4Info output
func (m MP4MetadataCommandProvider) ParseTitle(output string) (string, error) {
	return parseOutputStringForRegex(`(?msU)Title\s+(?P<content>\S.*$)`, output)
}

// ParseDescription parses the description from the MP4Info output
func (m MP4MetadataCommandProvider) ParseDescription(output string) (string, error) {
	return parseOutputStringForRegex(`(?msU)Comment\s+(?P<content>\S.*)Record date`, output)
}

// ParseCreator parses the creator from the MP4Info output
func (m MP4MetadataCommandProvider) ParseCreator(output string) (string, error) {
	return parseOutputStringForRegex(`(?msU)Artist\s+(?P<content>\S.*$)`, output)
}

// ParsePublishedAt parses the publishedAt from the MP4Info output
func (m MP4MetadataCommandProvider) ParsePublishedAt(output string) (*time.Time, error) {
	str, err := parseOutputStringForRegex(`(?msU)Record date\s+(?P<content>\S.*$)`, output)
	if err != nil {
		return nil, err
	}

	tm, err := time.Parse("2006-01-02", str)
	if err != nil {
		return nil, err
	}

	return &tm, nil
}

// ParseDuration parses the duration from the MP4Info output
func (m MP4MetadataCommandProvider) ParseDuration(output string) (*time.Duration, error) {
	hourRegex := `(?m)(?P<hour>\d*)\shr`
	minuteRegex := `(?m)(?P<minutes>\d*)\smin`
	secondsRegex := `(?m)(?P<seconds>\d*)\ss`

	hours, err := parseOutputStringForRegex(hourRegex, output)
	if err != nil {
		return nil, err
	}
	minutes, err := parseOutputStringForRegex(minuteRegex, output)
	if err != nil {
		return nil, err
	}
	seconds, err := parseOutputStringForRegex(secondsRegex, output)
	if err != nil {
		return nil, err
	}

	durationString := fmt.Sprintf("%sh%sm%ss", hours, minutes, seconds)
	duration, err := time.ParseDuration(durationString)

	return &duration, nil
}

func parseOutputStringForRegex(regex string, output string) (string, error) {
	re := regexp.MustCompile(regex)
	matches := re.FindStringSubmatch(output)

	if matches == nil {
		return "", fmt.Errorf("Failed to find content for regex %s", regex)
	}

	return strings.TrimSpace(matches[1]), nil
}

// Set sets metadata on an mp4 item
func (m MP4MetadataCommandProvider) Set(path string, metadata *videometadata.Metadata) error {
	values, err := buildTagEditorSetString(metadata)
	if err != nil {
		return err
	}

	return writeTagMetadata(values, path, &utils.OSCommand{})
}

func buildTagEditorSetString(metadata *videometadata.Metadata) (*[]string, error) {
	var valueString []string
	if metadata.Title != "" {
		valueString = append(valueString, fmt.Sprintf("title=%s", metadata.Title))
	}
	if metadata.Description != "" {
		valueString = append(valueString, fmt.Sprintf("comment=%s", metadata.Description))
	}
	if metadata.Creator != "" {
		valueString = append(valueString, fmt.Sprintf("artist=%s", metadata.Creator))
	}
	if metadata.PublishedAt != nil {
		publishString := metadata.PublishedAt.Format("2006-01-02")
		valueString = append(valueString, fmt.Sprintf("recorddate=%s", publishString))
	}

	if len(valueString) == 0 {
		return nil, errors.New("Provided metadata did not contain any data to write")
	}
	return &valueString, nil
}

func writeTagMetadata(value *[]string, path string, osCmd utils.OSCommandProvider) error {
	command := []string{
		"set",
		"-f",
		path,
		"--values",
	}

	command = append(command, (*value)...)

	out, err := osCmd.Run("tageditor", command...)
	if err != nil {
		return fmt.Errorf("Could not write metadata for %s.\nResponse from tageditor was: %s\nError %s", path, out, err)
	}

	return nil
}
