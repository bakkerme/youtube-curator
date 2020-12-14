package videometadata

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type testMetadataCommandProvider struct {
	Title       string
	Description string
	Creator     string
	PublishedAt *time.Time
	Duration    *time.Duration
}

func (m testMetadataCommandProvider) Run(path string) (string, error) {
	return "", nil
}

func (m testMetadataCommandProvider) Set(path string, metadata *Metadata) error {
	return nil
}

func (m testMetadataCommandProvider) ParseTitle(output string) (string, error) {
	if m.Title == "" {
		return "", errors.New("Bad Data")
	}
	return m.Title, nil
}

func (m testMetadataCommandProvider) ParseDescription(output string) (string, error) {
	if m.Description == "" {
		return "", errors.New("Bad Data")
	}
	return m.Description, nil
}

func (m testMetadataCommandProvider) ParseCreator(output string) (string, error) {
	if m.Creator == "" {
		return "", errors.New("Bad Data")
	}
	return m.Creator, nil
}

func (m testMetadataCommandProvider) ParsePublishedAt(output string) (*time.Time, error) {
	if m.PublishedAt == nil {
		return nil, errors.New("Bad Data")
	}
	return m.PublishedAt, nil
}

func (m testMetadataCommandProvider) ParseDuration(output string) (*time.Duration, error) {
	if m.Duration == nil {
		return nil, errors.New("Bad Data")
	}
	return m.Duration, nil
}

func testBrokenFieldsForMetadataParser(t *testing.T, field string, videoExpect *Metadata, metadataProvider *testMetadataCommandProvider) {
	metadata, pErr := parseVideoMetadataOutput("", metadataProvider)

	if pErr == nil || len(*pErr) != 1 {
		t.Errorf("Metadata Parser did not return error as expected")
	}

	if (*pErr)[0].field != field {
		t.Errorf("Field is incorrect. Expected %s, got %s", field, (*pErr)[0].field)
	}

	if !reflect.DeepEqual(*videoExpect, *metadata) {
		t.Errorf("Parsed output for %s did not match expected. Expected\n%+v\ngot\n%+v\n", field, videoExpect, metadata)
	}
}

func TestParseErrorResponse(t *testing.T) {
	expectedUnparsedFields := []string{
		"field1",
		"field2",
		"field3",
	}

	metadataParseError := handleParseErrorResponse(
		&[]fieldParseError{
			fieldParseError{
				errors.New("error 1"),
				"field1",
			},
			fieldParseError{
				errors.New("error 2"),
				"field2",
			},
			fieldParseError{
				errors.New("error 3"),
				"field3",
			},
		},
		"/ima/path",
	)

	if metadataParseError.err == "" {
		t.Error("metadataParseError should have an error string that's not empty")
	}

	if !reflect.DeepEqual(expectedUnparsedFields, metadataParseError.unparsedFields) {
		t.Errorf("Parsed output did not match expected. Expected\n%+v\ngot\n%+v\n", expectedUnparsedFields, metadataParseError.unparsedFields)
	}
}

func TestAppendParseError(t *testing.T) {
	t.Run("appendParseError doesn't append a nil fieldParseError", func(t *testing.T) {
		var fpe []fieldParseError
		appendParseError(&fpe, nil)

		if fpe != nil {
			t.Errorf("appendParseError appended something when nil was provided. Got %+v", fpe)
		}
	})

	t.Run("appendParseError appends a fieldParseError", func(t *testing.T) {
		var fpe []fieldParseError
		fpe = *appendParseError(&fpe, &fieldParseError{
			errors.New("error 1"),
			"field1",
		})

		if fpe == nil {
			t.Error("appendParseError didn't append a value")
		}

		if len(fpe) != 1 {
			t.Errorf("The fieldParseError slice should only contain 1 value. Got %+v", fpe)
		}
	})

	t.Run("appendParseError appends multiple fieldParseErrors", func(t *testing.T) {
		var fpe []fieldParseError
		fpe = *appendParseError(&fpe, &fieldParseError{
			errors.New("error 1"),
			"field1",
		})

		fpe = *appendParseError(&fpe, &fieldParseError{
			errors.New("error 2"),
			"field2",
		})

		if fpe == nil {
			t.Error("appendParseError didn't append a value")
		}

		if len(fpe) != 2 {
			t.Errorf("The fieldParseError slice should only contain 1 value. Got %+v", fpe)
		}
	})
}

func TestBuildParseError(t *testing.T) {
	t.Run("buildParseError return fieldParseError struct when provided error", func(t *testing.T) {
		field := "field1"
		fpe := buildParseError(field, errors.New("Bad Error"))

		if fpe == nil {
			t.Error("fieldParseError should not return nil")
		}

		if fpe.field != field {
			t.Errorf("fieldParseError field was not the field passed in. Expected %s, got %s", field, fpe.field)
		}
	})

	t.Run("buildParseError returns nil when not provided error", func(t *testing.T) {
		field := "field1"
		fpe := buildParseError(field, nil)

		if fpe != nil {
			t.Errorf("fieldParseError should return nil. Got %+v", fpe)
		}
	})
}

func buildVideoExpect() (*Metadata, error) {
	publishedAt, err := time.Parse("2006-01-02", "1992-05-01")
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration("3m")
	if err != nil {
		return nil, err
	}

	return &Metadata{
		Title:       "a title",
		Description: "a Description",
		Creator:     "a Creator",
		PublishedAt: &publishedAt,
		Duration:    &duration,
	}, nil
}

func TestMetaDataParser(t *testing.T) {
	videoExpect, err := buildVideoExpect()
	if err != nil {
		t.Error(err)
	}

	t.Run("A valid metadata file results in the correct metadata", func(t *testing.T) {
		metadataProvider := testMetadataCommandProvider{
			videoExpect.Title,
			videoExpect.Description,
			videoExpect.Creator,
			videoExpect.PublishedAt,
			videoExpect.Duration,
		}
		metadata, pErr := parseVideoMetadataOutput("", metadataProvider)
		if pErr != nil {
			t.Errorf("Video Metadata Parser return parse errors %+v", pErr)
		}

		if !reflect.DeepEqual(*metadata, *videoExpect) {
			t.Errorf("Parsed output did not match expected. Expected\n%+v\ngot\n%+v\n", *videoExpect, *metadata)
		}
	})

	t.Run("An invalid Title returns a VideoMetadataError with the Title field in the UnparsedFields", func(t *testing.T) {
		videoExpectWithoutTitle := *videoExpect
		videoExpectWithoutTitle.Title = ""

		metadataProvider := testMetadataCommandProvider{
			"",
			videoExpect.Description,
			videoExpect.Creator,
			videoExpect.PublishedAt,
			videoExpect.Duration,
		}

		testBrokenFieldsForMetadataParser(t, "Title", &videoExpectWithoutTitle, &metadataProvider)
	})

	t.Run("An invalid Description returns a VideoMetadataError with the Description field in the UnparsedFields", func(t *testing.T) {
		videoExpectWithoutDescription := *videoExpect
		videoExpectWithoutDescription.Description = ""

		metadataProvider := testMetadataCommandProvider{
			videoExpect.Title,
			"",
			videoExpect.Creator,
			videoExpect.PublishedAt,
			videoExpect.Duration,
		}

		testBrokenFieldsForMetadataParser(t, "Description", &videoExpectWithoutDescription, &metadataProvider)
	})

	t.Run("An invalid Creator returns a VideoMetadataError with the Creator field in the UnparsedFields", func(t *testing.T) {
		videoExpectWithoutCreator := *videoExpect
		videoExpectWithoutCreator.Creator = ""

		metadataProvider := testMetadataCommandProvider{
			videoExpect.Title,
			videoExpect.Description,
			"",
			videoExpect.PublishedAt,
			videoExpect.Duration,
		}

		testBrokenFieldsForMetadataParser(t, "Creator", &videoExpectWithoutCreator, &metadataProvider)
	})

	t.Run("An invalid PublishedAt returns a VideoMetadataError with the PublishedAt field in the UnparsedFields", func(t *testing.T) {
		videoExpectWithoutPublishedAt := *videoExpect
		videoExpectWithoutPublishedAt.PublishedAt = nil

		metadataProvider := testMetadataCommandProvider{
			videoExpect.Title,
			videoExpect.Description,
			videoExpect.Creator,
			nil,
			videoExpect.Duration,
		}

		testBrokenFieldsForMetadataParser(t, "PublishedAt", &videoExpectWithoutPublishedAt, &metadataProvider)
	})

	t.Run("An invalid Duration returns a VideoMetadataError with the Duration field in the UnparsedFields", func(t *testing.T) {
		videoExpectWithoutDuration := *videoExpect
		videoExpectWithoutDuration.Duration = nil

		metadataProvider := testMetadataCommandProvider{
			videoExpect.Title,
			videoExpect.Description,
			videoExpect.Creator,
			videoExpect.PublishedAt,
			nil,
		}

		testBrokenFieldsForMetadataParser(t, "Duration", &videoExpectWithoutDuration, &metadataProvider)
	})
}

func testBrokenFieldsForInfoResponse(t *testing.T, field string, videoExpect *Metadata, metadataProvider *testMetadataCommandProvider) {
	infoResponse := buildVideoMetadataResponse("", "/path", metadataProvider)

	metadata := infoResponse.Metadata
	fpErr := infoResponse.ParseError

	if fpErr.err == "" {
		t.Errorf("Video Info Error should have returned error string")
	}

	if fpErr.unparsedFields[0] != field {
		t.Errorf("Invalid file did not return unparsed field %s", field)
	}

	if !reflect.DeepEqual(*videoExpect, *metadata) {
		t.Errorf("Parsed output for %s did not match expected. Expected\n%+v\ngot\n%+v\n", field, videoExpect, metadata)
	}
}

func TestResponse(t *testing.T) {
	videoExpect, err := buildVideoExpect()
	if err != nil {
		t.Error(err)
	}

	t.Run("A valid metadata file results in correct metadata response", func(t *testing.T) {
		metadataProvider := testMetadataCommandProvider{
			videoExpect.Title,
			videoExpect.Description,
			videoExpect.Creator,
			videoExpect.PublishedAt,
			videoExpect.Duration,
		}
		infoResponse := buildVideoMetadataResponse("", "/path", &metadataProvider)

		if infoResponse.ParseError != nil {
			t.Errorf("Video Response returned parse error when none was expected. Error: %s", infoResponse.ParseError.Error())
		}

		if !reflect.DeepEqual(*infoResponse.Metadata, *videoExpect) {
			t.Errorf("Parsed output did not match expected. Expected\n%+v\ngot\n%+v\n", *videoExpect, *infoResponse.Metadata)
		}
	})

	t.Run("An invalid Title return a fieldParseError and other expected data", func(t *testing.T) {
		videoExpectWithoutTitle := *videoExpect
		videoExpectWithoutTitle.Title = ""

		metadataProvider := testMetadataCommandProvider{
			"",
			videoExpect.Description,
			videoExpect.Creator,
			videoExpect.PublishedAt,
			videoExpect.Duration,
		}

		testBrokenFieldsForInfoResponse(t, "Title", &videoExpectWithoutTitle, &metadataProvider)
	})

	t.Run("An invalid Description return a fieldParseError and other expected data", func(t *testing.T) {
		videoExpectWithoutDescription := *videoExpect
		videoExpectWithoutDescription.Description = ""

		metadataProvider := testMetadataCommandProvider{
			videoExpect.Title,
			"",
			videoExpect.Creator,
			videoExpect.PublishedAt,
			videoExpect.Duration,
		}

		testBrokenFieldsForInfoResponse(t, "Description", &videoExpectWithoutDescription, &metadataProvider)
	})

	t.Run("An invalid Creator return a fieldParseError and other expected data", func(t *testing.T) {
		videoExpectWithoutCreator := *videoExpect
		videoExpectWithoutCreator.Creator = ""

		metadataProvider := testMetadataCommandProvider{
			videoExpect.Title,
			videoExpect.Description,
			"",
			videoExpect.PublishedAt,
			videoExpect.Duration,
		}

		testBrokenFieldsForInfoResponse(t, "Creator", &videoExpectWithoutCreator, &metadataProvider)
	})

	t.Run("An invalid PublishedAt return a fieldParseError and other expected data", func(t *testing.T) {
		videoExpectWithoutPublishedAt := *videoExpect
		videoExpectWithoutPublishedAt.PublishedAt = nil

		metadataProvider := testMetadataCommandProvider{
			videoExpect.Title,
			videoExpect.Description,
			videoExpect.Creator,
			nil,
			videoExpect.Duration,
		}

		testBrokenFieldsForInfoResponse(t, "PublishedAt", &videoExpectWithoutPublishedAt, &metadataProvider)
	})

	t.Run("An invalid Duration return a fieldParseError and other expected data", func(t *testing.T) {
		videoExpectWithoutDuration := *videoExpect
		videoExpectWithoutDuration.Duration = nil

		metadataProvider := testMetadataCommandProvider{
			videoExpect.Title,
			videoExpect.Description,
			videoExpect.Creator,
			videoExpect.PublishedAt,
			nil,
		}

		testBrokenFieldsForInfoResponse(t, "Duration", &videoExpectWithoutDuration, &metadataProvider)
	})
}
