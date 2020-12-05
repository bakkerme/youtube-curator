package tageditor

import (
	"fmt"
	"testing"
	"time"
)

func TestStringParsers(t *testing.T) {
	commandProvider := MP4MetadataCommandProvider{}

	t.Run("ParseTitle should parse out a title", func(t *testing.T) {
		expectedTitle := "A Title"
		value := fmt.Sprintf("Title             %s", expectedTitle)
		title, err := commandProvider.ParseTitle(value)

		if err != nil {
			t.Error("ParseTitle should not return an error")
		}

		if title != expectedTitle {
			t.Errorf("ParseTitle did not return expected value. Expected %s, got %s", expectedTitle, title)
		}
	})

	t.Run("ParseTitle should return an error if there is no title", func(t *testing.T) {
		value := fmt.Sprintf("SOME BAD DATA")
		_, err := commandProvider.ParseTitle(value)

		if err == nil {
			t.Errorf("ParseTitle did not return an error on a bad input")
		}
	})

	t.Run("parseDescription should parse out a title", func(t *testing.T) {
		expectedDescription := `I'm a multiline
		description.`
		descriptionTemplate := `Comment           %s

			Record date`

		value := fmt.Sprintf(descriptionTemplate, expectedDescription)
		desc, err := commandProvider.ParseDescription(value)

		if err != nil {
			t.Errorf("parseDescription should not return an error. Error was %s", err)
		}

		if desc != expectedDescription {
			t.Errorf("parseDescription did not return expected value. \nExpected\n'%s'\n got\n'%s'", expectedDescription, desc)
		}
	})

	t.Run("parseDescription should return an error if there is no description", func(t *testing.T) {
		value := fmt.Sprintf("SOME BAD DATA")
		_, err := commandProvider.ParseDescription(value)

		if err == nil {
			t.Errorf("parseDescription did not return an error on a bad input")
		}
	})

	t.Run("parseCreator should parse out a creator", func(t *testing.T) {
		expectedCreator := "A Creator"
		creatorTemplate := "Artist           %s"

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
		value := fmt.Sprintf("SOME BAD DATA")
		_, err := commandProvider.ParseCreator(value)

		if err == nil {
			t.Errorf("parseCreator did not return an error on a bad input")
		}
	})

	t.Run("parsePublishedAt should parse out a publishedAt", func(t *testing.T) {
		expectedPublishedAt := "A PublishedAt"
		publishedAtTemplate := `Record date           %s`

		value := fmt.Sprintf(publishedAtTemplate, expectedPublishedAt)
		publishedAt, err := commandProvider.ParsePublishedAt(value)

		if err != nil {
			t.Error("parsePublishedAt should not return an error")
		}

		if publishedAt != expectedPublishedAt {
			t.Errorf("parsePublishedAt did not return expected value. Expected %s, got %s", expectedPublishedAt, publishedAt)
		}
	})

	t.Run("parsePublishedAt should return an error if there is no publishedAt", func(t *testing.T) {
		value := fmt.Sprintf("SOME BAD DATA")
		_, err := commandProvider.ParsePublishedAt(value)

		if err == nil {
			t.Errorf("parsePublishedAt did not return an error on a bad input")
		}
	})

	t.Run("parseDuration should parse out a duration", func(t *testing.T) {
		var expectedDuration time.Duration
		expectedDuration = 4555000000000

		value := "Duration                      1 hr 15 min 55 s 942 ms 290 µs 200 ns"
		duration, err := commandProvider.ParseDuration(value)

		if err != nil {
			t.Error("parseDuration should not return an error")
		}

		if duration != expectedDuration {
			t.Errorf("parseDuration did not return expected value. Expected %d, got %d", expectedDuration, duration)
		}
	})

	t.Run("parseDuration should return an error if there is no duration", func(t *testing.T) {
		value := fmt.Sprintf("SOME BAD DATA")
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
