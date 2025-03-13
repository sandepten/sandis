package main

import "strings"

func parseInput(input string) string {
	data := strings.Split(input, " ")[0]
	if strings.Contains(data, "*") || strings.Contains(data, "$") {
		return ""
	}
	return data
}
