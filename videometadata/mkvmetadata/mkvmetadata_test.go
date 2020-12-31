package mkvmetadata

import (
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/videometadata"
	"reflect"
	"testing"
	"time"
)

type readMockOSCommand struct {
	expectedCommand *string
	expectedParams  *[]string
	output          *string
	returnError     bool
	t               *testing.T
	callCount       int
}

func (osc *readMockOSCommand) Run(name string, arg ...string) (*[]byte, error) {
	if osc.returnError {
		return nil, errors.New("Gee willikers, looks like the server did a flopsy boingo")
	}

	if osc.expectedCommand != nil && *osc.expectedCommand != name {
		osc.t.Errorf("Run recieved wrong command. Expected %+v, got %+v", osc.expectedCommand, name)
	}

	if osc.expectedParams != nil && !reflect.DeepEqual(*osc.expectedParams, arg) {
		osc.t.Errorf("Run recieved wrong params. Expected %+v, got %+v", osc.expectedParams, arg)
	}

	if osc.output != nil {
		out := []byte(*osc.output)
		osc.callCount++
		return &out, nil
	}

	osc.callCount++
	out := []byte("")
	return &out, nil
}

func TestReadMetadata(t *testing.T) {
	t.Run("getMetadata should run two commands with correct results", func(t *testing.T) {
		path := "/path"
		commands := "mkvinfo"
		params := []string{path}

		output := "output1"

		result, err := getMetadata(path, &readMockOSCommand{
			expectedCommand: &commands,
			expectedParams:  &params,
			output:          &output,
			t:               t,
		})
		if err != nil {
			t.Error(err)
		}

		if result != output {
			t.Errorf("getMetadata did not return correct output. Expected %s, got %s", output, result)
		}
	})

	t.Run("getMetadata should fail", func(t *testing.T) {
		_, err := getMetadata("/path", &readMockOSCommand{returnError: true, t: t})
		if err == nil {
			t.Error("Should have returned error")
		}
	})
}

func TestStringParsers(t *testing.T) {
	commandProvider := MKVMetadataCommandProvider{}

	t.Run("ParseTitle should parse out a title", func(t *testing.T) {
		expectedTitle := "A Title"
		value := fmt.Sprintf("| + Title: %s", expectedTitle)
		title, err := commandProvider.ParseTitle(value)

		if err != nil {
			t.Error("ParseTitle should not return an error")
		}

		if title != expectedTitle {
			t.Errorf("ParseTitle did not return expected value. Expected %s, got %s", expectedTitle, title)
		}
	})

	t.Run("ParseTitle should return an error if there is no title", func(t *testing.T) {
		value := fmt.Sprintf("| SOME BAD DATA")
		_, err := commandProvider.ParseTitle(value)

		if err == nil {
			t.Errorf("ParseTitle did not return an error on a bad input")
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
		desc, err := commandProvider.ParseDescription(value)

		if err != nil {
			t.Error("parseDescription should not return an error")
		}

		if desc != expectedDescription {
			t.Errorf("parseDescription did not return expected value. Expected %s, got %s", expectedDescription, desc)
		}
	})

	t.Run("parseDescription should return an error if there is no description", func(t *testing.T) {
		value := fmt.Sprintf("| SOME BAD DATA")
		_, err := commandProvider.ParseDescription(value)

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
		creator, err := commandProvider.ParseCreator(value)

		if err != nil {
			t.Error("parseCreator should not return an error")
		}

		if creator != expectedCreator {
			t.Errorf("parseCreator did not return expected value. Expected %s, got %s", expectedCreator, creator)
		}
	})

	t.Run("parseCreator should return an error if there is no creator", func(t *testing.T) {
		value := fmt.Sprintf("| SOME BAD DATA")
		_, err := commandProvider.ParseCreator(value)

		if err == nil {
			t.Errorf("parseCreator did not return an error on a bad input")
		}
	})

	t.Run("parsePublishedAt should parse out a publishedAt", func(t *testing.T) {
		expectedPublishedAtString := "20201117"
		rawInputString := `| + Name: DATE
|  + String: %s
|`

		value := fmt.Sprintf(rawInputString, expectedPublishedAtString)

		expectedPublishedAt, err := time.Parse("20060102", expectedPublishedAtString)
		if err != nil {
			t.Error(err)
		}

		publishedAt, err := commandProvider.ParsePublishedAt(value)

		if err != nil {
			t.Errorf("parsePublishedAt should not return an error %s", err)
		}

		if *publishedAt != expectedPublishedAt {
			t.Errorf("parsePublishedAt did not return expected value. Expected %s, got %s", expectedPublishedAt, publishedAt)
		}
	})

	t.Run("parsePublishedAt should return an error if there is no publishedAt", func(t *testing.T) {
		value := fmt.Sprintf("| SOME BAD DATA")
		_, err := commandProvider.ParsePublishedAt(value)

		if err == nil {
			t.Errorf("parsePublishedAt did not return an error on a bad input")
		}
	})

	t.Run("parsePublishedAt should return an error if the publishedAt date is invalid", func(t *testing.T) {
		value := `| + Name: DATE
|  + String: asdfl
|`
		_, err := commandProvider.ParsePublishedAt(value)

		if err == nil {
			t.Errorf("parsePublishedAt did not return an error on a bad input")
		}
	})

	t.Run("parseDuration should parse out a duration", func(t *testing.T) {
		var expectedDuration time.Duration
		expectedDuration = 1692000000000

		value := "| + Duration: 00:28:12.563000000"
		duration, err := commandProvider.ParseDuration(value)

		if err != nil {
			t.Error("parseDuration should not return an error")
		}

		if *duration != expectedDuration {
			t.Errorf("parseDuration did not return expected value. Expected %s, got %s", expectedDuration, duration)
		}
	})

	t.Run("parseDuration should return an error if there is no duration", func(t *testing.T) {
		value := fmt.Sprintf("| SOME BAD DATA")
		_, err := commandProvider.ParseDuration(value)

		if err == nil {
			t.Errorf("parseDuration did not return an error on a bad input")
		}
	})

	t.Run("parseDuration should return an error if the duration is invalid", func(t *testing.T) {
		value := "| + Duration: asdfsdfsdf"
		_, err := commandProvider.ParseDuration(value)

		if err == nil {
			t.Errorf("parseDuration did not return an error on a bad input")
		}
	})

	t.Run("parseDuration should parse out a duration", func(t *testing.T) {
		value := "| + Duration: 00:28:-20.563000000"
		_, err := commandProvider.ParseDuration(value)

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

func TestCommands(t *testing.T) {
	t.Run("MKVMetadataCommandProvider should throw error for Set since it's unimpemented :)", func(t *testing.T) {
		var cmdProv videometadata.CommandProvider
		cmdProv = MKVMetadataCommandProvider{}

		err := cmdProv.Set("", nil)
		if err == nil {
			t.Error("Who implemented my function without telling me!!!! wow")
		}
	})
}

func TestInterface(t *testing.T) {
	t.Run("MKVMetadataCommandProvider should be a CommandProvider", func(t *testing.T) {
		var cmdProv videometadata.CommandProvider
		cmdProv = MKVMetadataCommandProvider{}

		t.Logf("lets use cmdProv so the compiler doesn't get mad %s", cmdProv)
	})
}
