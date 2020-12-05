package tageditor

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// MP4MetadataCommandProvider provides the ability for the VideoMetadata
type MP4MetadataCommandProvider struct{}

// Run MP4Info on the provided file
func (m MP4MetadataCommandProvider) Run(path string) (string, error) {
	out, err := exec.Command("./tageditor-3.3.10.AppImage", "get", "-f", path).Output()
	if err != nil {
		return "", fmt.Errorf("Could not load metadata for %s.\nResponse from mkvinfo was: %s\nError %s", path, out, err)
	}

	infoOut, err := exec.Command("./tageditor-3.3.10.AppImage", "info", "-f", path).Output()
	if err != nil {
		return "", fmt.Errorf("Could not load metadata for %s.\nResponse from mkvinfo was: %s\nError %s", path, out, err)
	}

	return string(out) + string(infoOut), nil
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
func (m MP4MetadataCommandProvider) ParsePublishedAt(output string) (string, error) {
	return parseOutputStringForRegex(`(?msU)Record date\s+(?P<content>\S.*$)`, output)
}

// ParseDuration parses the duration from the MP4Info output
func (m MP4MetadataCommandProvider) ParseDuration(output string) (time.Duration, error) {
	hourRegex := `(?m)(?P<hour>\d*)\shr`
	minuteRegex := `(?m)(?P<hour>\d*)\smin`
	secondsRegex := `(?m)(?P<seconds>\d*)\ss`

	hours, err := parseOutputStringForRegex(hourRegex, output)
	if err != nil {
		return 0, err
	}
	minutes, err := parseOutputStringForRegex(minuteRegex, output)
	if err != nil {
		return 0, err
	}
	seconds, err := parseOutputStringForRegex(secondsRegex, output)
	if err != nil {
		return 0, err
	}

	durationString := fmt.Sprintf("%sh%sm%ss", hours, minutes, seconds)
	duration, err := time.ParseDuration(durationString)

	return duration, nil
}

func parseOutputStringForRegex(regex string, output string) (string, error) {
	re := regexp.MustCompile(regex)
	matches := re.FindStringSubmatch(output)

	if matches == nil {
		return "", fmt.Errorf("Failed to find content for regex %s", regex)
	}

	return strings.TrimSpace(matches[1]), nil
}
