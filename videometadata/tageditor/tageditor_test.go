package tageditor

import (
	"errors"
	"fmt"
	"hyperfocus.systems/youtube-curator-server/videometadata"
	"reflect"
	"testing"
	"time"
)

type readMockOSCommand struct {
	expectedCommand *[]string
	expectedParams  *[][]string
	output          *[]string
	returnError     *[]bool
	t               *testing.T
	callCount       int
}

func (osc *readMockOSCommand) Run(name string, arg ...string) (*[]byte, error) {
	if osc.returnError != nil && (*osc.returnError)[osc.callCount] {
		return nil, errors.New("Gee willikers, looks like the server did a flopsy boingo")
	}

	if osc.expectedCommand != nil && (*osc.expectedCommand)[osc.callCount] != name {
		osc.t.Errorf("Run recieved wrong command. Expected %+v, got %+v", (*osc.expectedCommand)[osc.callCount], name)
	}

	if osc.expectedParams != nil && !reflect.DeepEqual((*osc.expectedParams)[osc.callCount], arg) {
		osc.t.Errorf("Run recieved wrong params. Expected %+v, got %+v", (*osc.expectedParams)[osc.callCount], arg)
	}

	if osc.output != nil {
		out := []byte((*osc.output)[osc.callCount])
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
		commands := &[]string{"tageditor", "tageditor"}
		params := &[][]string{[]string{"get", "-f", path}, []string{"info", "-f", path}}

		output1 := "output1"
		output2 := "output2"
		output := &[]string{output1, output2}

		result, err := getMetadata(path, &readMockOSCommand{expectedCommand: commands, expectedParams: params, output: output, t: t})
		if err != nil {
			t.Error(err)
		}

		if result != output1+output2 {
			t.Errorf("getMetadata did not return correct output. Expected %s, got %s", output1+output2, result)
		}
	})

	t.Run("getMetadata should fail the first command", func(t *testing.T) {
		_, err := getMetadata("/path", &readMockOSCommand{returnError: &[]bool{true, false}})
		if err == nil {
			t.Error("Should have returned error")
		}
	})

	t.Run("getMetadata should fail the second command", func(t *testing.T) {
		_, err := getMetadata("/path", &readMockOSCommand{returnError: &[]bool{false, true}})
		if err == nil {
			t.Error("Should have returned error")
		}
	})
}

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
		expectedPublishedAtString := "2017-06-15"
		expectedPublishedAt, err := time.Parse("2006-01-02", expectedPublishedAtString)
		if err != nil {
			t.Errorf("time.Parse returned an error %s", err)
		}

		publishedAtTemplate := `Record date           %s`

		value := fmt.Sprintf(publishedAtTemplate, expectedPublishedAtString)
		publishedAt, err := commandProvider.ParsePublishedAt(value)

		if err != nil {
			t.Errorf("parsePublishedAt should not return an error %s", err)
		}

		if *publishedAt != expectedPublishedAt {
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

		if *duration != expectedDuration {
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

type writeMockOSCommand struct {
	expectedCommand string
	expectedParams  *[]string
	returnError     bool
	t               *testing.T
}

func (osc *writeMockOSCommand) Run(name string, arg ...string) (*[]byte, error) {
	if osc.returnError {
		return nil, errors.New("Gee willikers, looks like the server did a flopsy boingo")
	}

	if osc.expectedCommand != name {
		osc.t.Errorf("Run recieved wrong command. Expected %s, got %s", osc.expectedCommand, name)
	}

	if !reflect.DeepEqual(*osc.expectedParams, arg) {
		osc.t.Errorf("Run recieved wrong params. Expected %s, got %s", osc.expectedParams, arg)
	}

	return nil, nil
}

func TestWrite(t *testing.T) {
	t.Run("buildTagEditorSetString returns a valid string for a metadata input", func(t *testing.T) {
		publishedAt, err := time.Parse("2006-01-02", "1992-05-01")
		if err != nil {
			t.Error("Failed to parse time string")
		}

		mt := &videometadata.Metadata{
			Title:       "a title",
			Description: "a description",
			Creator:     "a creator",
			PublishedAt: &publishedAt,
		}

		str, err := buildTagEditorSetString(mt)
		if err != nil {
			t.Errorf("buildTagEditorSetString returned an error when none was expected. Got %s", err)
		}

		expected := []string{
			"title=" + mt.Title,
			"comment=" + mt.Description,
			"artist=" + mt.Creator,
			"recorddate=1992-05-01",
		}

		if !reflect.DeepEqual(*str, expected) {
			t.Errorf("buildTagEditorSetString did not return correct result. Expected\n%s\ngot\n%s", expected, str)
		}
	})

	t.Run("buildTagEditorSetString should return an error if no valid tags are to be written", func(t *testing.T) {
		mt := &videometadata.Metadata{
			Title:       "",
			Description: "",
			Creator:     "",
			PublishedAt: nil,
		}

		_, err := buildTagEditorSetString(mt)
		if err == nil {
			t.Error("buildTagEditorSetString did not return an error")
		}
	})

	t.Run("writeTagMetadata runs the correct command provided", func(t *testing.T) {
		path := "/path"
		value := "title=test"
		value2 := "title2=test2"
		values := &[]string{value, value2}
		mosc := &writeMockOSCommand{
			expectedCommand: "tageditor",
			expectedParams:  &[]string{"set", "-f", path, "--values", value, value2},
			t:               t,
		}

		err := writeTagMetadata(values, path, mosc)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("writeTagMetadata returns error if command fails", func(t *testing.T) {
		mosc := &writeMockOSCommand{
			t:           t,
			returnError: true,
		}

		err := writeTagMetadata(&[]string{"string1", "string2"}, "/path", mosc)
		if err == nil {
			t.Error("writeTagMetadata should return error if command returns error")
		}
	})

}

func TestInterface(t *testing.T) {
	t.Run("MP4MetadataCommandProvider should be a CommandProvider", func(t *testing.T) {
		var cmdProv videometadata.CommandProvider
		cmdProv = MP4MetadataCommandProvider{}

		t.Logf("lets use cmdProv so the compiler doesn't get mad %s", cmdProv)
	})
}
