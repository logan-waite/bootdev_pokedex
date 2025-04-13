package main

import (
	"strings"
)

func cleanInput(str string) []string {
	result := []string{}
	list := strings.Fields(str)

	for _, item := range list {
		item = strings.ToLower(item)
		result = append(result, item)
	}

	return result
}
