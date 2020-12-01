package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

var mkvExpect = mkvInfoMetadata{
	title: "1970's Camera Tech: How they showed you what settings to use",
	description: `Ever wonder how pro cameras from the 1970's worked? Learn about their single most important tool for the photographer (and lots else!) in this exposé.
Strings of text which take you places!`,
	creator:     "Technology Connections",
	publishedAt: "20201030",
	duration:    1692000000000,
}

func testBrokenFieldsForMetadataParser(t *testing.T, mkvExpect *mkvInfoMetadata, field string) {
	fileName := fmt.Sprintf("./testfiles/mkv-broken-%s.txt", field)
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Errorf("Could not load %s off disk. %s", fileName, err)
	}

	metadata, pErr := parseMKVInfoOutput(string(file))

	if pErr == nil || len(*pErr) != 1 {
		t.Errorf("MKV Metadata Parser did not return error as expected")
	}

	if (*pErr)[0].field != field {
		t.Errorf("Invalid file did not return unparsed field %s", field)
	}

	if !reflect.DeepEqual(*mkvExpect, *metadata) {
		t.Errorf("Parsed output for %s did not match expected. Expected\n%+v\ngot\n%+v\n", field, mkvExpect, metadata)
	}
}

func testBrokenFieldsForInfoResponse(t *testing.T, mkvExpect *mkvInfoMetadata, field string) {
	fileName := fmt.Sprintf("./testfiles/mkv-broken-%s.txt", field)
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Errorf("Could not load %s off disk. %s", fileName, err)
	}

	infoResponse := buildMKVInfoResponse(string(file), fileName)

	metadata := infoResponse.metadata
	fpErr := infoResponse.parseError

	if fpErr.err == "" {
		t.Errorf("MKV Info Error should have returned error string")
	}

	if fpErr.unparsedFields[0] != field {
		t.Errorf("Invalid file did not return unparsed field %s", field)
	}

	if !reflect.DeepEqual(*mkvExpect, *metadata) {
		t.Errorf("Parsed output for %s did not match expected. Expected\n%+v\ngot\n%+v\n", field, mkvExpect, metadata)
	}
}

func TestMetaDataParser(t *testing.T) {
	t.Run("A valid metadata file results in the correct metadata", func(t *testing.T) {
		fileName := "./testfiles/mkv.txt"
		file, err := ioutil.ReadFile(fileName)
		if err != nil {
			t.Errorf("Could not load %s off disk. %s", fileName, err)
		}

		metadata, pErr := parseMKVInfoOutput(string(file))
		if pErr != nil {
			t.Errorf("MKV Metadata Parser failed with error %s", err.Error())
		}

		if !reflect.DeepEqual(*metadata, mkvExpect) {
			t.Errorf("Parsed output did not match expected. Expected\n%+v\ngot\n%+v\n", mkvExpect, metadata)
		}
	})

	t.Run("An invalid Title returns a VideoMetadataError with the Title field in the UnparsedFields", func(t *testing.T) {
		mkvExpectWithoutTitle := mkvExpect
		mkvExpectWithoutTitle.title = ""

		testBrokenFieldsForMetadataParser(t, &mkvExpectWithoutTitle, "title")
	})

	t.Run("An invalid Description returns a VideoMetadataError with the Description field in the UnparsedFields", func(t *testing.T) {
		mkvExpectWithoutDescription := mkvExpect
		mkvExpectWithoutDescription.description = ""

		testBrokenFieldsForMetadataParser(t, &mkvExpectWithoutDescription, "description")
	})

	t.Run("An invalid Creator returns a VideoMetadataError with the Creator field in the UnparsedFields", func(t *testing.T) {
		mkvExpectWithoutCreator := mkvExpect
		mkvExpectWithoutCreator.creator = ""

		testBrokenFieldsForMetadataParser(t, &mkvExpectWithoutCreator, "creator")
	})

	t.Run("An invalid PublishedAt returns a VideoMetadataError with the PublishedAt field in the UnparsedFields", func(t *testing.T) {
		mkvExpectWithoutPublishedAt := mkvExpect
		mkvExpectWithoutPublishedAt.publishedAt = ""

		testBrokenFieldsForMetadataParser(t, &mkvExpectWithoutPublishedAt, "publishedAt")
	})

	t.Run("An invalid Duration returns a VideoMetadataError with the Duration field in the UnparsedFields", func(t *testing.T) {
		mkvExpectWithoutDuration := mkvExpect
		mkvExpectWithoutDuration.duration = 0

		testBrokenFieldsForMetadataParser(t, &mkvExpectWithoutDuration, "duration")
	})
}

func TestResponse(t *testing.T) {
	t.Run("A valid metadata file results in correct metadata response", func(t *testing.T) {
		fileName := "./testfiles/mkv.txt"
		file, err := ioutil.ReadFile(fileName)
		if err != nil {
			t.Errorf("Could not load %s off disk. %s", fileName, err)
		}

		infoResponse := buildMKVInfoResponse(string(file), fileName)

		if infoResponse.parseError != nil {
			t.Errorf("MKV Response returned parse error when none was expected. Error: %s", err.Error())
		}

		if !reflect.DeepEqual(*infoResponse.metadata, mkvExpect) {
			t.Errorf("Parsed output did not match expected. Expected\n%+v\ngot\n%+v\n", mkvExpect, infoResponse.metadata)
		}
	})

	t.Run("An invalid Title return a fieldParseError and other expected data", func(t *testing.T) {
		mkvExpectWithoutTitle := mkvExpect
		mkvExpectWithoutTitle.title = ""

		testBrokenFieldsForInfoResponse(t, &mkvExpectWithoutTitle, "title")
	})

	t.Run("An invalid Description return a fieldParseError and other expected data", func(t *testing.T) {
		mkvExpectWithoutDescription := mkvExpect
		mkvExpectWithoutDescription.description = ""

		testBrokenFieldsForInfoResponse(t, &mkvExpectWithoutDescription, "description")
	})

	t.Run("An invalid Creator return a fieldParseError and other expected data", func(t *testing.T) {
		mkvExpectWithoutCreator := mkvExpect
		mkvExpectWithoutCreator.creator = ""

		testBrokenFieldsForInfoResponse(t, &mkvExpectWithoutCreator, "creator")
	})

	t.Run("An invalid PublishedAt return a fieldParseError and other expected data", func(t *testing.T) {
		mkvExpectWithoutPublishedAt := mkvExpect
		mkvExpectWithoutPublishedAt.publishedAt = ""

		testBrokenFieldsForInfoResponse(t, &mkvExpectWithoutPublishedAt, "publishedAt")
	})

	t.Run("An invalid Duration return a fieldParseError and other expected data", func(t *testing.T) {
		mkvExpectWithoutDuration := mkvExpect
		mkvExpectWithoutDuration.duration = 0

		testBrokenFieldsForInfoResponse(t, &mkvExpectWithoutDuration, "duration")
	})
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

func TestStringParsers(t *testing.T) {
	t.Run("parseTitle should parse out a title", func(t *testing.T) {
		expectedTitle := "A Title"
		value := fmt.Sprintf("| + Title: %s", expectedTitle)
		title, err := parseTitle(value)

		if err != nil {
			t.Error("parseTitle should not return an error")
		}

		if title != expectedTitle {
			t.Errorf("parseTitle did not return expected value. Expected %s, got %s", expectedTitle, title)
		}
	})

	t.Run("parseTitle should return an error if there is no title", func(t *testing.T) {
		value := fmt.Sprintf("| SOME BAD DATA")
		_, err := parseTitle(value)

		if err == nil {
			t.Errorf("parseTitle did not return an error on a bad input")
		}
	})

	t.Run("parseDescription should parse out a title", func(t *testing.T) {
		expectedDescription := `I'm a multiline
		description.
		`
		descriptionTemplate := `| + Name: DESCRIPTION
|  + String: %s
|`

		value := fmt.Sprintf(descriptionTemplate, expectedDescription)
		desc, err := parseDescription(value)

		if err != nil {
			t.Error("parseDescription should not return an error")
		}

		if desc != expectedDescription {
			t.Errorf("parseDescription did not return expected value. Expected %s, got %s", expectedDescription, desc)
		}
	})

	t.Run("parseDescription should return an error if there is no description", func(t *testing.T) {
		value := fmt.Sprintf("| SOME BAD DATA")
		_, err := parseDescription(value)

		if err == nil {
			t.Errorf("parseDescription did not return an error on a bad input")
		}
	})

	t.Run("parseCreator should parse out a creator", func(t *testing.T) {
		expectedCreator := "A Creator"
		creatorTemplate := `| + Name: ARTIST
|  + String: %s
|`

		value := fmt.Sprintf(creatorTemplate, expectedCreator)
		creator, err := parseCreator(value)

		if err != nil {
			t.Error("parseCreator should not return an error")
		}

		if creator != expectedCreator {
			t.Errorf("parseCreator did not return expected value. Expected %s, got %s", expectedCreator, creator)
		}
	})

	t.Run("parseCreator should return an error if there is no creator", func(t *testing.T) {
		value := fmt.Sprintf("| SOME BAD DATA")
		_, err := parseCreator(value)

		if err == nil {
			t.Errorf("parseCreator did not return an error on a bad input")
		}
	})

	t.Run("parsePublishedAt should parse out a publishedAt", func(t *testing.T) {
		expectedPublishedAt := "A PublishedAt"
		publishedAt := `| + Name: DATE
|  + String: %s
|`

		value := fmt.Sprintf(publishedAt, expectedPublishedAt)
		publishedAt, err := parsePublishedAt(value)

		if err != nil {
			t.Error("parsePublishedAt should not return an error")
		}

		if publishedAt != expectedPublishedAt {
			t.Errorf("parsePublishedAt did not return expected value. Expected %s, got %s", expectedPublishedAt, publishedAt)
		}
	})

	t.Run("parsePublishedAt should return an error if there is no publishedAt", func(t *testing.T) {
		value := fmt.Sprintf("| SOME BAD DATA")
		_, err := parsePublishedAt(value)

		if err == nil {
			t.Errorf("parsePublishedAt did not return an error on a bad input")
		}
	})

	t.Run("parseDuration should parse out a duration", func(t *testing.T) {
		var expectedDuration time.Duration
		expectedDuration = 1692000000000

		value := "| + Duration: 00:28:12.563000000"
		duration, err := parseDuration(value)

		if err != nil {
			t.Error("parseDuration should not return an error")
		}

		if duration != expectedDuration {
			t.Errorf("parseDuration did not return expected value. Expected %s, got %s", expectedDuration, duration)
		}
	})

	t.Run("parseDuration should return an error if there is no duration", func(t *testing.T) {
		value := fmt.Sprintf("| SOME BAD DATA")
		_, err := parseDuration(value)

		if err == nil {
			t.Errorf("parseDuration did not return an error on a bad input")
		}
	})

	t.Run("parseOutputStringForRegex should return correct output for regex", func(t *testing.T) {
		expected := "asdf"
		out, err := parseOutputStringForRegex(`(?msU)Test: (?P<content>.*$)`, fmt.Sprintf("Test: %s", expected))

		if err != nil {
			t.Error("parseOutputStringForRegex should not return an error")
		}

		if out != expected {
			t.Errorf("parseOutputStringForRegex should have returned ")
		}
	})

	t.Run("parseOutputStringForRegex should return error when not matching regex", func(t *testing.T) {
		_, err := parseOutputStringForRegex(`(?msU)Test: (?P<content>.*$)`, "BAD RESULT")

		if err == nil {
			t.Error("parseOutputStringForRegex should return an error")
		}
	})
}
