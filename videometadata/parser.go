package videometadata

import (
	"fmt"
)

// fieldParseError is a parse error for a single field
type fieldParseError struct {
	err   error
	field string
}

func buildVideoMetadataResponse(output string, path string, pr CommandProvider) *Response {
	metadata, pErrs := parseVideoMetadataOutput(output, pr)
	if pErrs != nil {
		return &Response{
			metadata,
			handleParseErrorResponse(pErrs, path),
		}
	}

	return &Response{
		metadata,
		nil,
	}
}

func parseVideoMetadataOutput(output string, pr CommandProvider) (*Metadata, *[]fieldParseError) {
	var parseErrors []fieldParseError

	title, err := pr.ParseTitle(output)
	parseErrors = *appendParseError(&parseErrors, buildParseError("Title", err))

	description, err := pr.ParseDescription(output)
	parseErrors = *appendParseError(&parseErrors, buildParseError("Description", err))

	creator, err := pr.ParseCreator(output)
	parseErrors = *appendParseError(&parseErrors, buildParseError("Creator", err))

	publishedAt, err := pr.ParsePublishedAt(output)
	parseErrors = *appendParseError(&parseErrors, buildParseError("PublishedAt", err))

	duration, err := pr.ParseDuration(output)
	parseErrors = *appendParseError(&parseErrors, buildParseError("Duration", err))

	if parseErrors != nil && len(parseErrors) > 0 {
		return &Metadata{
			Title:       title,
			Description: description,
			Creator:     creator,
			PublishedAt: publishedAt,
			Duration:    duration,
		}, &parseErrors
	}

	return &Metadata{
		Title:       title,
		Description: description,
		Creator:     creator,
		PublishedAt: publishedAt,
		Duration:    duration,
	}, nil

}

func handleParseErrorResponse(pErrs *[]fieldParseError, videoPath string) *ParseError {
	errString := fmt.Sprintf("Could not load metadata for %s. Failed on parse stage", videoPath)
	var unparsedFields []string
	for _, pErr := range *pErrs {
		errString += fmt.Sprintf("\n - %s", pErr.err.Error())
		unparsedFields = append(unparsedFields, pErr.field)
	}

	return &ParseError{
		errString,
		unparsedFields,
	}
}

func appendParseError(parseErrors *[]fieldParseError, pErr *fieldParseError) *[]fieldParseError {
	if pErr != nil {
		appendedParseErrors := append(*parseErrors, *pErr)
		parseErrors = &appendedParseErrors
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
