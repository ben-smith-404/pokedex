package main

import (
	"strings"
)

func cleanInput(text string) []string {
	var strs []string
	split_strs := strings.Split(text, " ")
	for _, str := range split_strs {
		if len(str) > 0 {
			strs = append(strs, strings.ToLower((str)))
		}
	}
	return strs
}
