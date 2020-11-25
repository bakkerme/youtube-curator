package main

import (
	"regexp"
)

func getDescriptionFromOutput(output string) string {
	re := regexp.MustCompile(`(?msU)\+ Name: DESCRIPTION.*String: (.*)^\|`)
	return re.FindAllString(output, -1)[0]
}
