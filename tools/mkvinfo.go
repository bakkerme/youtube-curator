package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type mkvInfoResponse struct {
	metadata   *mkvInfoMetadata
	parseError *mkvMetadataParseError
}

type mkvInfoMetadata struct {
	title       string
	description string
	creator     string
	publishedAt string
	duration    time.Duration
}

// mkvMetadataParseError implements VideoMetadataError
type mkvMetadataParseError struct {
	err            string
	unparsedFields []string
}

func (pErr *mkvMetadataParseError) Error() string {
	return pErr.err
}

func (pErr *mkvMetadataParseError) UnparsedFields() []string {
	return pErr.unparsedFields
}

// fieldParseError is a parse error for a single field
type fieldParseError struct {
	err   error
	field string
}

func loadVideoFile(path string) (*mkvInfoResponse, error) {
	cmd := exec.Command(fmt.Sprintf("mkvinfo \"%s\"", path))
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Could not load metadata for %s. Error %s", path, err)
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Could not load metadata for %s. Failed on output stage. Error %s", path, err)
	}

	return buildMKVInfoResponse(string(out), path), nil
}

func buildMKVInfoResponse(output string, path string) *mkvInfoResponse {
	metadata, pErrs := parseMKVInfoOutput(output)
	if pErrs != nil {
		return &mkvInfoResponse{
			metadata,
			handleParseErrorResponse(pErrs, path),
		}
	}

	return &mkvInfoResponse{
		metadata,
		nil,
	}
}

func parseMKVInfoOutput(output string) (*mkvInfoMetadata, *[]fieldParseError) {
	var parseErrors []fieldParseError

	title, err := parseTitle(output)
	parseErrors = *appendParseError(&parseErrors, buildParseError("title", err))

	description, err := parseDescription(output)
	parseErrors = *appendParseError(&parseErrors, buildParseError("description", err))

	creator, err := parseCreator(output)
	parseErrors = *appendParseError(&parseErrors, buildParseError("creator", err))

	publishedAt, err := parsePublishedAt(output)
	parseErrors = *appendParseError(&parseErrors, buildParseError("publishedAt", err))

	duration, err := parseDuration(output)
	parseErrors = *appendParseError(&parseErrors, buildParseError("duration", err))

	if parseErrors != nil && len(parseErrors) > 0 {
		return &mkvInfoMetadata{
			title:       title,
			description: description,
			creator:     creator,
			publishedAt: publishedAt,
			duration:    duration,
		}, &parseErrors
	}

	return &mkvInfoMetadata{
		title:       title,
		description: description,
		creator:     creator,
		publishedAt: publishedAt,
		duration:    duration,
	}, nil

}

func handleParseErrorResponse(pErrs *[]fieldParseError, videoPath string) *mkvMetadataParseError {
	errString := fmt.Sprintf("Could not load metadata for %s. Failed on parse stage", videoPath)
	var unparsedFields []string
	for _, pErr := range *pErrs {
		errString += fmt.Sprintf("\n - %s", pErr.err.Error())
		unparsedFields = append(unparsedFields, pErr.field)
	}

	return &mkvMetadataParseError{
		errString,
		unparsedFields,
	}
}

func appendParseError(parseErrors *[]fieldParseError, pErr *fieldParseError) *[]fieldParseError {
	if pErr != nil {
		fmt.Println(parseErrors)
		appendedParseErrors := append(*parseErrors, *pErr)
		parseErrors = &appendedParseErrors
		fmt.Println(parseErrors)
	}

	return parseErrors
}

func buildParseError(field string, err error) *fieldParseError {
	if err != nil {
		return &fieldParseError{
			fmt.Errorf("Field %s: Error %s", field, err),
			field,
		}
	}

	return nil
}

func parseTitle(output string) (string, error) {
	return parseOutputStringForRegex(`(?m)\| \+ Title: (?P<content>.*$)`, output)
}

func parseDescription(output string) (string, error) {
	return parseOutputStringForRegex(`(?msU)\+ Name: DESCRIPTION.*String: (?P<content>.*)\n^\|`, output)
}

func parseCreator(output string) (string, error) {
	return parseOutputStringForRegex(`(?msU)ARTIST.*String: (?P<content>.*$)`, output)
}

func parsePublishedAt(output string) (string, error) {
	return parseOutputStringForRegex(`(?msU)DATE.*String: (?P<content>.*$)`, output)
}

func parseDuration(output string) (time.Duration, error) {
	str, err := parseOutputStringForRegex(`(?msU)Duration: (?P<content>.*$)`, output)
	if err != nil {
		return 0, err
	}

	milisecondless := strings.Split(str, ".")[0]
	units := strings.Split(milisecondless, ":")
	durationString := fmt.Sprintf("%sh%sm%ss", units[0], units[1], units[2])
	duration, err := time.ParseDuration(durationString)

	if err != nil {
		return 0, err
	}

	return duration, nil
}

func parseOutputStringForRegex(regex string, output string) (string, error) {
	re := regexp.MustCompile(regex)
	matches := re.FindStringSubmatch(output)

	if matches == nil {
		return "", fmt.Errorf("Failed to find content for regex %s", regex)
	}

	return matches[1], nil
}
